# acodex（Codex 多账号快速切换）

<p align="center">
  <img alt="acodex" src="https://img.shields.io/badge/acodex-Codex%20%E5%A4%9A%E8%B4%A6%E5%8F%B7%E5%88%87%E6%8D%A2-blue" />
  <a href="./LICENSE"><img alt="License" src="https://img.shields.io/badge/License-MIT-green.svg" /></a>
  <img alt="Platform" src="https://img.shields.io/badge/platform-Windows-0078d4" />
  <img alt="Go" src="https://img.shields.io/badge/Go-1.24%2B-00ADD8" />
</p>

<p align="center">
<strong>中文</strong> | <a href="./README.en.md">English</a>
</p>

`acodex` 是一个 **Windows 命令行工具**，用来保存和切换 Codex 的登录凭据文件 `auth.json`，从而在多个账号之间快速来回切换。

- 默认 Codex 认证文件：`%USERPROFILE%\\.codex\\auth.json`
- 支持 `CODEX_HOME`：认证文件为 `%CODEX_HOME%\\auth.json`

## 功能特性

- 多 profile 管理：`save/use/list/current/delete`
- 切换时安全保护：
    - 若当前 `auth.json` **不属于任何已保存 profile**，会提示确认是否覆盖
    - 覆盖前自动备份：`auth.json.bak.<unix_ts>`
- 一键打开数据目录：`acodex open`
- 首次运行自动安装：自动复制到 `~\\.acodex\\bin\\acodex.exe` 并写入当前用户 PATH

## 运行环境

- Windows（目前依赖注册表 PATH 和 `explorer`）
- Go（从源码构建时）版本以 `go.mod` 为准

## 安装

### 方式 A：从 GitHub Releases 下载（推荐）

1. 在 Releases 下载 `acodex.exe`
2. 双击或在终端运行一次：它会自动安装到：`%USERPROFILE%\\.acodex\\bin\\acodex.exe`
3. **关闭并重新打开** 终端（cmd/PowerShell），然后就可以直接用 `acodex ...`

> 说明：第一次运行会把 `~\\.acodex\\bin` 追加到 **当前用户** 的 `Path`（注册表 `HKCU\\Environment`）。

### 方式 B：从源码构建

```powershell
go build -o acodex.exe .
```

## 快速开始

下面以两个账号 `work` / `personal` 为例：

```powershell
acodex save work
# 用 Codex 登录另一个账号，生成新的 auth.json
acodex save personal

acodex use work
acodex use personal
```

> 注意：`save` 是 **移动**（rename）当前 `auth.json` 到 profile 目录，执行后原位置会没有 `auth.json`。

## 命令说明

- `acodex save <alias>`：保存当前 `auth.json` 到新 profile（移动）
- `acodex use <alias>`：切换到指定 profile（复制到 Codex 目录，必要时确认 + 自动备份）
- `acodex list`：列出 profiles（`*` 标记当前）
- `acodex current`：输出当前 alias
- `acodex delete <alias>`：删除 profile
- `acodex open`：用资源管理器打开 `%USERPROFILE%\\.acodex`

## 数据目录

- Codex 凭据：
    - 默认：`%USERPROFILE%\\.codex\\auth.json`
    - 可选：`%CODEX_HOME%\\auth.json`
- acodex：
    - Profiles：`%USERPROFILE%\\.acodex\\profiles\\<alias>\\auth.json`
    - 当前标记：`%USERPROFILE%\\.acodex\\current`
    - 程序安装：`%USERPROFILE%\\.acodex\\bin\\acodex.exe`

## License

MIT
