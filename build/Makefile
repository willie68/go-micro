## Taken from http://www.codershaven.com/multi-platform-makefile-for-go/
## Could be a good possibility if more than one platform must be supported
## It also use the git version as version for the binary.
EXECUTABLE=bin/gomicro-service
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN=$(EXECUTABLE)_darwin_amd64
VERSION=$(shell git describe --tags --always --long --dirty)
LDFLAGS=-s -w -X main.version=$(VERSION)

windows: $(WINDOWS) ## Build for Windows

linux: $(LINUX) ## Build for Linux

darwin: $(DARWIN) ## Build for Darwin (macOS)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build -i -v -o $(WINDOWS) -ldflags="$(LDFLAGS)"  

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -i -v -o $(LINUX) -ldflags="$(LDFLAGS)"

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -i -v -o $(DARWIN) -ldflags="$(LDFLAGS)"


build: windows linux darwin ## Build binaries
	@echo version: $(VERSION)

all: test build ## Build and run tests

test: ## Run unit tests
	go test

clean: ## Remove previous build
	rm -f $(WINDOWS) $(LINUX) $(DARWIN)

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'