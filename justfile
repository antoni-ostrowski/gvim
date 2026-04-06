default:
  just --list

dev:
  go run ./cmd/

build:
  go build -o gvim ./cmd 

debug:
  tail -f /tmp/gvim.log
