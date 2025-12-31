package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"
)

const (
	AppDirName = ".acodex"
	BinDirName = "bin"
	ExeName    = "acodex.exe"
)

var isChinese bool

func init() {
	isChinese = detectLanguage()
}

func detectLanguage() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Control Panel\International`, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer k.Close()

	val, _, err := k.GetStringValue("LocaleName")
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.ToLower(val), "zh")
}

func tr(en, zh string) string {
	if isChinese {
		return zh
	}
	return en
}

// Path helpers

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

func appHome() string {
	return filepath.Join(homeDir(), AppDirName)
}

func binDir() string {
	return filepath.Join(appHome(), BinDirName)
}

func installedExePath() string {
	return filepath.Join(binDir(), ExeName)
}

func codexHome() string {
	if v := os.Getenv("CODEX_HOME"); v != "" {
		return v
	}
	return filepath.Join(homeDir(), ".codex")
}

func authPath() string {
	return filepath.Join(codexHome(), "auth.json")
}

func profilesDir() string {
	return filepath.Join(appHome(), "profiles")
}

func currentFile() string {
	return filepath.Join(appHome(), "current")
}

// Auto-installation logic

func isInstalled() bool {
	self, err := os.Executable()
	if err != nil {
		return false
	}
	self, _ = filepath.Abs(self)
	target, _ := filepath.Abs(installedExePath())
	return strings.EqualFold(self, target)
}

func autoInstall() error {
	if err := os.MkdirAll(binDir(), 0755); err != nil {
		return err
	}

	self, err := os.Executable()
	if err != nil {
		return err
	}

	if !strings.EqualFold(self, installedExePath()) {
		if err := copyFile(self, installedExePath()); err != nil {
			return err
		}
	}

	return addToUserPath(binDir())
}

func addToUserPath(dir string) error {
	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		"Environment",
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	path, _, _ := key.GetStringValue("Path")
	parts := strings.Split(path, ";")
	for _, p := range parts {
		if strings.EqualFold(p, dir) {
			return nil
		}
	}

	if path == "" {
		path = dir
	} else {
		path = path + ";" + dir
	}
	return key.SetStringValue("Path", path)
}

// Utility functions

func ensureDirs() error {
	return os.MkdirAll(profilesDir(), 0755)
}

func sha256File(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func readCurrent() string {
	b, err := os.ReadFile(currentFile())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func writeCurrent(v string) {
	_ = os.WriteFile(currentFile(), []byte(v), 0644)
}

func listProfiles() ([]string, error) {
	entries, err := os.ReadDir(profilesDir())
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func findMatchAuth(p string) string {
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	h, _ := sha256File(p)
	profiles, _ := listProfiles()
	for _, a := range profiles {
		ap := filepath.Join(profilesDir(), a, "auth.json")
		if hh, _ := sha256File(ap); hh == h {
			return a
		}
	}
	return ""
}

func confirm(msg string) bool {
	fmt.Print(msg, " (y/N): ")
	in := bufio.NewReader(os.Stdin)
	s, _ := in.ReadString('\n')
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "y" || s == "yes"
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

func backupIfExists(p string) {
	if _, err := os.Stat(p); err == nil {
		_ = os.Rename(p, p+".bak."+fmt.Sprint(time.Now().Unix()))
	}
}

// Command implementations

func cmdSave(alias string) error {
	if alias == "" {
		return errors.New(tr("alias required", "需要指定别名"))
	}
	if err := ensureDirs(); err != nil {
		return err
	}
	src := authPath()
	if _, err := os.Stat(src); err != nil {
		return errors.New(tr("auth.json not found", "未找到 auth.json"))
	}
	dst := filepath.Join(profilesDir(), alias, "auth.json")
	if _, err := os.Stat(dst); err == nil {
		return errors.New(tr("profile exists", "配置文件已存在"))
	}
	_ = os.MkdirAll(filepath.Dir(dst), 0755)
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	writeCurrent(alias)
	fmt.Println(tr("saved:", "已保存:"), alias)
	return nil
}

func cmdUse(alias string) error {
	src := filepath.Join(profilesDir(), alias, "auth.json")
	if _, err := os.Stat(src); err != nil {
		return errors.New(tr("profile not found", "未找到配置文件"))
	}
	dst := authPath()
	if _, err := os.Stat(dst); err == nil {
		if findMatchAuth(dst) == "" && !confirm(tr("current auth.json will be overwritten", "当前 auth.json 将被覆盖")) {
			return errors.New(tr("aborted", "已取消"))
		}
		backupIfExists(dst)
	}
	if err := copyFile(src, dst); err != nil {
		return err
	}
	writeCurrent(alias)
	fmt.Println(tr("using:", "当前使用:"), alias)
	return nil
}

// Main entry point

func main() {
	if !isInstalled() {
		if err := autoInstall(); err != nil {
			fmt.Println(tr("install failed:", "安装失败:"), err)
			os.Exit(1)
		}
		fmt.Println(tr("acodex installed successfully. Please restart your terminal.", "acodex 已自动安装完成，请重新打开终端使用。"))
		return
	}

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "save":
		if len(os.Args) < 3 {
			fmt.Println(tr("usage: acodex save <alias>", "用法: acodex save <别名>"))
			return
		}
		if err := cmdSave(os.Args[2]); err != nil {
			fmt.Println(tr("error:", "错误:"), err)
		}
	case "use":
		if len(os.Args) < 3 {
			fmt.Println(tr("usage: acodex use <alias>", "用法: acodex use <别名>"))
			return
		}
		if err := cmdUse(os.Args[2]); err != nil {
			fmt.Println(tr("error:", "错误:"), err)
		}
	case "list":
		p, _ := listProfiles()
		cur := readCurrent()
		for _, a := range p {
			if a == cur {
				fmt.Println("*", a)
			} else {
				fmt.Println(" ", a)
			}
		}
	case "current":
		fmt.Println(readCurrent())
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println(tr("usage: acodex delete <alias>", "用法: acodex delete <别名>"))
			return
		}
		_ = os.RemoveAll(filepath.Join(profilesDir(), os.Args[2]))
		fmt.Println(tr("deleted:", "已删除:"), os.Args[2])
	case "open":
		_ = exec.Command("explorer", appHome()).Start()
	default:
		fmt.Println(tr("unknown command", "未知命令"))
		printUsage()
	}
}

func printUsage() {
	fmt.Println(tr("Usage: acodex <command> [args]", "用法: acodex <命令> [参数]"))
	fmt.Println(tr("\nAvailable commands:", "\n可用命令:"))

	cmds := []struct {
		name string
		en   string
		zh   string
	}{
		{"save <alias>", "Save current auth.json as a profile", "保存当前 auth.json 为配置文件"},
		{"use <alias>", "Switch to a specific profile", "切换到指定配置文件"},
		{"list", "List all saved profiles", "列出所有保存的配置文件"},
		{"current", "Show current profile name", "显示当前配置文件名称"},
		{"delete <alias>", "Delete a profile", "删除指定配置文件"},
		{"open", "Open application directory", "打开应用程序目录"},
	}

	for _, cmd := range cmds {
		fmt.Printf("  %-15s %s\n", cmd.name, tr(cmd.en, cmd.zh))
	}
}
