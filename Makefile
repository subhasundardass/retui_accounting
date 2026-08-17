.PHONY: help build debug-build debug kill attach run package clean tidy screen test lint log dev check-main-pkg

.DEFAULT_GOAL := build

# ──────────────────────────────────────────────────────────────────────────
# Colors (these are Make variables, not shell env vars — `${GREEN}` inside a
# recipe line is expanded by Make itself before the shell ever sees it,
# since Make accepts both $() and ${} for its own variable references).
# ──────────────────────────────────────────────────────────────────────────
GREEN  := \033[0;32m
RED    := \033[0;31m
NC     := \033[0m

# ──────────────────────────────────────────────────────────────────────────
# Entry point
#
# MAIN_PKG is auto-detected by scanning for main.go files, skipping
# vendor/build/.git and any codegen tool under cmd/*-gen. Override
# explicitly when you have more than one real entrypoint:
#   make build MAIN_PKG=./cmd/yourapp
# ──────────────────────────────────────────────────────────────────────────
BINARY_NAME := retui
BUILD_DIR   := build
CONFIG_SRC  := internal/config/config.yml
MAIN_PKG := ./cmd/app

ifndef MAIN_PKG
MAIN_PKG_CANDIDATES := $(shell find . \
	\( -path './vendor' -o -path './.git' -o -path './$(BUILD_DIR)' -o -name '*-gen' \) -prune -o \
	-name 'main.go' -print 2>/dev/null | sed -e 's|/main.go$$||' -e 's|^\./||' | sort -u)
NUM_CANDIDATES := $(words $(MAIN_PKG_CANDIDATES))
ifeq ($(NUM_CANDIDATES),1)
MAIN_PKG := ./$(MAIN_PKG_CANDIDATES)
endif
endif

GOOS    ?= $(shell go env GOOS)
GOARCH  ?= $(shell go env GOARCH)
LDFLAGS := -s -w

BIN := $(BINARY_NAME)
DEBUG_BIN := $(BINARY_NAME)-debug
ifeq ($(GOOS),windows)
	BIN := $(BINARY_NAME).exe
	DEBUG_BIN := $(BINARY_NAME)-debug.exe
endif

DLV_PORT ?= 2345

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

# Fails fast with a clear message instead of the old cryptic
# "directory not found" error, and tells you exactly what to type.
check-main-pkg:
	@if [ -z "$(strip $(MAIN_PKG))" ]; then \
		printf "${RED}==> ERROR: could not determine MAIN_PKG automatically.${NC}\n"; \
		if [ -n "$(strip $(MAIN_PKG_CANDIDATES))" ]; then \
			echo "==> Found multiple candidate main packages:"; \
			for p in $(MAIN_PKG_CANDIDATES); do echo "       - ./$$p"; done; \
			echo "==> Pick one: make build MAIN_PKG=./path/to/package"; \
		else \
			echo "==> No main.go found under this directory."; \
			echo "==> Pass it explicitly: make build MAIN_PKG=./cmd/yourapp"; \
		fi; \
		exit 1; \
	fi

## build: compile the release binary and assemble build/{binary,config.yml,data/}
build: check-main-pkg
	@printf "${GREEN}==> Building $(BINARY_NAME) ($(GOOS)/$(GOARCH)) from $(MAIN_PKG)${NC}\n"
	@mkdir -p $(BUILD_DIR)/data
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BIN) $(MAIN_PKG)
	@if [ -f "$(CONFIG_SRC)" ]; then \
		cp "$(CONFIG_SRC)" "$(BUILD_DIR)/config.yml"; \
	else \
		printf "${RED}==> WARNING: $(CONFIG_SRC) not found, skipping config copy${NC}\n"; \
	fi
	@echo "==> Build complete: $(BUILD_DIR)/$(BIN)"
	@ls -la $(BUILD_DIR)

## debug-build: compile an unstripped, unoptimized binary for use with dlv
##              (release builds use -ldflags "-s -w", which strips the
##              symbols a debugger needs — this is a separate artifact)
debug-build: check-main-pkg
	@printf "${GREEN}==> Building debug binary from $(MAIN_PKG)${NC}\n"
	@mkdir -p $(BUILD_DIR)/data
	@go build -gcflags "all=-N -l" -o $(BUILD_DIR)/$(DEBUG_BIN) $(MAIN_PKG)
	@if [ -f "$(CONFIG_SRC)" ]; then cp "$(CONFIG_SRC)" "$(BUILD_DIR)/config.yml"; fi
	@echo "==> Debug build complete: $(BUILD_DIR)/$(DEBUG_BIN)"

## run: build, then execute from inside build/ so relative paths resolve
##      the same way they would for an end user launching the binary directly
run: build
	@echo "==> Running $(BUILD_DIR)/$(BIN)"
	@cd $(BUILD_DIR) && ./$(BIN)

## package: build then zip build/ into a single distributable archive
package: build
	@cd $(BUILD_DIR) && zip -r ../$(BINARY_NAME)-$(GOOS)-$(GOARCH).zip .
	@echo "==> Packaged: $(BINARY_NAME)-$(GOOS)-$(GOARCH).zip"

## kill: free DLV_PORT if a previous debug session was left running
kill:
	@lsof -ti :$(DLV_PORT) | xargs -r kill -9 2>/dev/null || true
	@echo "==> Port $(DLV_PORT) is now free."

## debug: rebuild a debug binary and start Delve headless against it
debug: kill debug-build
	dlv exec $(BUILD_DIR)/$(DEBUG_BIN) \
		--headless \
		--listen=:$(DLV_PORT) \
		--api-version=2 \
		--accept-multiclient

## attach: start Delve against the last debug-build without rebuilding
##         (use this if you haven't changed any .go files since debug-build)
attach: kill
	@if [ ! -f "$(BUILD_DIR)/$(DEBUG_BIN)" ]; then \
		printf "${RED}==> No debug binary found — run 'make debug-build' first.${NC}\n"; \
		exit 1; \
	fi
	dlv exec $(BUILD_DIR)/$(DEBUG_BIN) \
		--headless \
		--listen=:$(DLV_PORT) \
		--api-version=2 \
		--accept-multiclient

## dev: run the app directly via `go run`, no build artifact produced
dev: check-main-pkg
	go run $(MAIN_PKG) --debug

## clean: remove the build directory
clean:
	@rm -rf $(BUILD_DIR)
	@echo "==> Cleaned $(BUILD_DIR)"

## tidy: clean up go.mod and go.sum
tidy:
	@printf "${GREEN}==> Tidying modules...${NC}\n"
	go mod tidy
	go mod verify

## screen: generate a new screen (usage: make screen NAME=general)
screen:
	@if [ -z "$(NAME)" ]; then \
		printf "${RED}==> Please provide a screen name: make screen NAME=general${NC}\n"; \
		exit 1; \
	fi
	@printf "${GREEN}==> Generating screen: $(NAME)...${NC}\n"
	@go run ./cmd/retui-gen -make screen:$(NAME)

## test: run the project's test workflow script
test:
	./test-workflow.sh

## lint: run go vet and golangci-lint if available
lint:
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || \
		echo "==> golangci-lint not installed, skipped (go vet still ran)"

## log: tail the application log
log:
	tail -n +1 -f retui.log
