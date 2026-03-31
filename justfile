default:
  just --list

dev:
  go run ./cmd/

build:
  go build ./cmd 

debug:
  tail -f /tmp/gvim.log
