SHELL=/bin/bash -o pipefail

ARTIFACT_NAME := terraform-provider-ethereum

GOPROXY ?= ""

ARTIFACTS_DIR ?= _artifacts

all: build

_artifacts:
	mkdir -p ${ARTIFACTS_DIR}

_bin:
	mkdir -p bin

define compile
	$(eval os = $1)
	$(eval extension = $2)
	@echo "building $(os) binary"
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build -o bin/$(ARTIFACT_NAME)_$(os)_amd64$(extension)
endef

.PHONY: build
build: _bin lint
	$(call compile, darwin)
	$(call compile, linux)
	$(call compile, windows, .exe)

.PHONY: test
test: _artifacts lint go_mod
	@echo "Missing tests"

.PHONY: lint
lint:
	@echo "Executing linters"

	golangci-lint run

.PHONY: go_mod
go_mod:
	go mod verify
