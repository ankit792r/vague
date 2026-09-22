# Vague — build and install
#
#   make              # build bin/vague (frontend + Go binary)
#   make install      # install to ~/.local (see PREFIX below)
#   make uninstall
#   make test
#
# Install elsewhere:
#   make install PREFIX=/usr/local
#   sudo make install PREFIX=/usr/local

ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

GO   ?= go
BUN  ?= bun

PREFIX ?= $(HOME)/.local
BINDIR ?= $(PREFIX)/bin

BIN := $(ROOT)/bin/vague

.PHONY: all build frontend go install uninstall test clean help

.DEFAULT_GOAL := build

help:
	@echo "Vague Makefile targets:"
	@echo "  make build       Build frontend and $(BIN)"
	@echo "  make install     Install binary, desktop file, icon, user systemd unit"
	@echo "  make uninstall   Remove install artifacts"
	@echo "  make test        Run Go tests"
	@echo "  make clean       Remove build outputs"
	@echo ""
	@echo "Variables: PREFIX=$(PREFIX)  BINDIR=$(BINDIR)  GO=$(GO)  BUN=$(BUN)"

all: build

build:
	$(ROOT)/scripts/build.sh

frontend:
	cd "$(ROOT)/frontend" && $(BUN) install && $(BUN) run build

go: frontend
	mkdir -p "$(ROOT)/bin"
	CGO_ENABLED=0 $(GO) build -ldflags "-s -w" -trimpath -o "$(BIN)" "$(ROOT)"

install:
	PREFIX="$(PREFIX)" BINDIR="$(BINDIR)" "$(ROOT)/scripts/install.sh" install

uninstall:
	PREFIX="$(PREFIX)" BINDIR="$(BINDIR)" "$(ROOT)/scripts/install.sh" uninstall

test:
	cd "$(ROOT)" && $(GO) test ./...

clean:
	rm -f "$(BIN)"
	rm -rf "$(ROOT)/bonefire/webview/output"
	rm -rf "$(ROOT)/frontend/node_modules"
