# gvim

A terminal-based text editor written in Go, inspired by Vim.



https://github.com/user-attachments/assets/3daa5812-1b9f-45c9-a9f1-c48e7933b584




> [!NOTE]
> Its really early :)

# Editor commands

- `:w` - write to opened file (open via `gvim ./some-file.txt`)
- `:q` - quit (or CTRL+C)

# Implemented motions
- `o` & `O` (line inserts)
- `$` & `0` (line jumps)
- `hjkl`
- arrow navigation

# Installation

Download binary from [latest release](https://github.com/antoni-ostrowski/gvim/releases)

# Run

```bash
./gvim ./some-file.txt
```

# Development

#### Prerequisites

- Go 1.25.4
- [just - task runner](https://github.com/casey/just)

```bash
just build
```

#### Debug Logs

Project uses a debug log utility for development:

```go
utils.Debuglog("cursor x = %v, cursor y = %v", x, y)
```

```bash
# read current debug output
just debug

```
