SRC=$(shell find . -name "*.go")
BIN="./bin"

.PHONY: vet deps clean

build: ensure-dir build-linux build-windows build-darwin compress

all: vet build	

deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

clean:
	$(info ******************** downloading dependencies ********************)
	rm -rf $(BIN)

ensure-dir:
	$(info ******************** ensure dir ********************)
	rm -rf bin
	mkdir bin

build-linux:
	$(info ******************** build linux ********************)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tomba.linux-amd64 *.go

build-windows:
	$(info ******************** build windows ********************)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tomba.windows-amd64.exe *.go

build-darwin:
	$(info ******************** build darwin ********************)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/tomba.darwin-amd-64 *.go

compress:
	$(info ******************** compress ********************)
	cd ./bin && find . -name 'tomba*' | xargs -I{} tar czf {}.tar.gz {}

vet:
	$(info ******************** vetting ********************)
	go vet ./...