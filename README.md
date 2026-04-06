# gvim

A terminal-based text editor written in Go, inspired by Vim.

## Prerequisites

- Go 1.25.4

## Build

```bash
just build
```

## Run

```bash
./gvim
# or
go run ./cmd
```

## Development

### Debug Logs

The project uses a debug log utility for development:

```go
utils.Debuglog("cursor x = %v, cursor y = %v", x, y)
```

```bash
# read current debug output
just debug

```

