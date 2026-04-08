# gvim

A terminal-based text editor written in Go, inspired by Vim.

# Editor commands

- `:w` - write to opened file (open via `gvim ./some-file.txt`)
- `:q` - quit (or CTRL+C)

# Run

```bash
./gvim ./some-file.txt
```

# Development

## Build

## Prerequisites

- Go 1.25.4
- [just - task runner](https://github.com/casey/just)

```bash
just build
```

### Debug Logs

Project uses a debug log utility for development:

```go
utils.Debuglog("cursor x = %v, cursor y = %v", x, y)
```

```bash
# read current debug output
just debug

```
