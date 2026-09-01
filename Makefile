SRC=$(shell find . -name "*.go")
BIN="./bin"

VERSION ?= dev
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w \
	-X github.com/tomba-io/tomba/pkg/version.Version=$(VERSION) \
	-X github.com/tomba-io/tomba/pkg/version.BuildDate=$(BUILD_DATE) \
	-X github.com/tomba-io/tomba/pkg/version.Commit=$(COMMIT)

.PHONY: vet deps clean demos

build: ensure-dir build-linux build-windows build-darwin build-darwin-arm64 compress

all: vet build

deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

clean:
	$(info ******************** clean bin ********************)
	rm -rf $(BIN)

ensure-dir:
	$(info ******************** ensure dir ********************)
	rm -rf bin
	mkdir bin

build-linux:
	$(info ******************** build linux ********************)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tomba.linux-amd64 *.go

build-windows:
	$(info ******************** build windows ********************)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tomba.windows-amd64.exe *.go

build-darwin:
	$(info ******************** build darwin ********************)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tomba.darwin-amd64 *.go

build-darwin-arm64:
	$(info ******************** build darwin arm64 ********************)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/tomba.darwin-arm64 *.go

compress:
	$(info ******************** compress ********************)
	cd ./bin && find . -name 'tomba*' | xargs -I{} tar czf {}.tar.gz {}

demos:
	$(info ******************** generating demos ********************)
	./bin/generate-demos

vet:
	$(info ******************** vetting ********************)
	go vet ./...
