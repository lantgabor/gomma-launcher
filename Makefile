VERSION ?= $(shell git describe --tags --abbrev=0 2> /dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD)

# GTK is linked through cgo, so builds target the machine they run on.
# Cross-compiling (e.g. GOOS=linux on macOS) would need a C toolchain and
# GTK libraries for the target, so build each platform on that platform.
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)
PKG     := ./cmd/gomma-launcher
LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT)

.PHONY: help
help:  ## Show this help message
	@echo "Available make targets:"
	@grep -E "^[a-zA-Z0-9 -]+:.*##" $(MAKEFILE_LIST) | sort | column -t -s "##"

.PHONY: build
build:  ## Build gomma-launcher for this machine
	go build -ldflags "$(LDFLAGS)" -o ./build/gomma-launcher $(PKG)

.PHONY: run
run:  ## Run gomma-launcher
	go run -ldflags "$(LDFLAGS)" $(PKG)

.PHONY: clean
clean:  ## Delete artifacts
	rm -rf build
	rm -rf release

.PHONY: release
release: clean build  ## Generate artifacts for release
	mkdir -p release
	cd build && tar cvzf ../release/gomma-launcher-$(VERSION)-$(GOOS)-$(GOARCH).tar.gz gomma-launcher
	@echo
	@echo "Generated artifacts:"
	@find release -maxdepth 1 -type f
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"

.PHONY: install
install:  ## Install gomma-launcher
	go install -ldflags "$(LDFLAGS)" $(PKG)
