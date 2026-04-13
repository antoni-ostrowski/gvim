binary_name := "gvim"
version := "0.1.0"

default:
    just --list

dev:
    go run -race ./cmd/ ./testfile.txt

# Build optimized binary
build:
    go build -o {{binary_name}}  ./cmd/

# Install the binary to your $GOPATH/bin
install:
    go install ./cmd/

# Clean up binaries and log files
clean:
    rm -f {{binary_name}}
    rm -f /tmp/gvim.log

debug:
    touch /tmp/gvim.log
    tail -f /tmp/gvim.log

lint:
    go vet ./...
