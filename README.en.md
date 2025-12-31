# acodex (Codex account switcher)

<p align="center">
<strong>English</strong> | <a href="./README.md">中文</a>
</p>

A **Windows CLI** to save and switch Codex `auth.json` profiles, so you can quickly switch between multiple Codex
accounts.

<p align="center">
  <img alt="acodex" src="https://img.shields.io/badge/acodex-Codex%20profile%20switcher-blue" />
  <a href="./LICENSE"><img alt="License" src="https://img.shields.io/badge/License-MIT-green.svg" /></a>
  <img alt="Platform" src="https://img.shields.io/badge/platform-Windows-0078d4" />
  <img alt="Go" src="https://img.shields.io/badge/Go-1.24%2B-00ADD8" />
</p>

## What it does

Codex stores credentials in `auth.json`. This tool manages multiple copies (profiles) and switches them for you.

- Default Codex auth path: `%USERPROFILE%\\.codex\\auth.json`
- Or set `CODEX_HOME` (auth file: `%CODEX_HOME%\\auth.json`)

## Features

- Manage profiles: `save/use/list/current/delete`
- Safer switching:
    - If current `auth.json` doesn’t match any saved profile, it asks before overwriting
    - Automatic backup before overwrite: `auth.json.bak.<unix_ts>`
- `acodex open`: open app data folder
- One-time self-install: copies itself to `%USERPROFILE%\\.acodex\\bin\\acodex.exe` and adds it to the **user PATH**

## Installation

### Option A: Download from GitHub Releases (recommended)

1. Download `acodex.exe` from Releases
2. Run it once. It will self-install to:
    - `%USERPROFILE%\\.acodex\\bin\\acodex.exe`
3. Re-open your terminal (cmd/PowerShell), then you can use `acodex ...`

### Option B: Build from source

```powershell
go build -o acodex.exe .
```

## Quick start

Example with `work` / `personal`:

```powershell
acodex save work
# login with Codex to generate a new auth.json
acodex save personal

acodex use work
acodex use personal
```

> Note: `save` **moves** (renames) `auth.json` into the profile folder.

## Commands

- `acodex save <alias>`: move current `auth.json` into a new profile
- `acodex use <alias>`: copy profile `auth.json` to Codex location (with confirmation + backup)
- `acodex list`: list profiles (`*` marks current)
- `acodex current`: print current alias
- `acodex delete <alias>`: delete a profile
- `acodex open`: open `%USERPROFILE%\\.acodex` in Explorer

## Data locations

- Profiles: `%USERPROFILE%\\.acodex\\profiles\\<alias>\\auth.json`
- Current marker: `%USERPROFILE%\\.acodex\\current`
- Installed binary: `%USERPROFILE%\\.acodex\\bin\\acodex.exe`

## License

MIT
