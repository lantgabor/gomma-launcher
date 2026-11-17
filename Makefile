VERSION ?= $(shell git describe --tags --abbrev=0 2> /dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD)

.PHONY: help
help:  ## Show this help message
	@echo "Available make targets:"
	@grep -E "^[a-zA-Z0-9 -]+:.*##" $(MAKEFILE_LIST) | sort | column -t -s "##"

.PHONY: build-linux
build-linux:  ## Build gomma-launcher for linux
	GOOS=linux GOARCH=amd64 go build \
		-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" \
		-o ./build/gomma-launcher \
		cmd/gomma-launcher/main.go

.PHONY: build
build: build-linux  ## Build gomma-launcher

.PHONY: run
run:  ## Run gomma-launcher
	GOOS=linux GOARCH=amd64 go run \
		-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" \
		./cmd/gomma-launcher/main.go

.PHONY: clean
clean:  ## Delete artifacts
	rm -rf build
	rm -rf release

.PHONY: release
release: clean build ## Generate artifacts for release
	mkdir -p release
	cd build && tar cvzf ../release/gomma-launcher-$(VERSION)-linux-amd64.tar.gz gomma-launcher
	@echo
	@echo "Generated artifacts:"
	@find release -maxdepth 1 -type f
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"

.PHONY: install
install:  ## Install gomma-launcher
	go install\
		-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}"
