.PHONY: help test test-conformance build lint update-spec clean

help:
	@echo "Available targets:"
	@echo "  make test              - Run all tests (auto-initializes submodule)"
	@echo "  make test-conformance  - Run conformance tests only"
	@echo "  make build             - Build all packages"
	@echo "  make lint              - Run linters"
	@echo "  make update-spec       - Update spec submodule to latest version"
	@echo "  make clean             - Clean build artifacts"
	@echo "  make help              - Show this help message"

test:
	@if [ ! -d testdata/spec/schema ]; then \
		echo "Initializing submodule..."; \
		git submodule update --init --recursive; \
	fi
	go test ./...

test-conformance:
	@if [ ! -d testdata/spec/schema ]; then \
		echo "Initializing submodule..."; \
		git submodule update --init --recursive; \
	fi
	go test -run TestConformance ./...

build:
	go build ./...

lint:
	go vet ./...

update-spec:
	git submodule update --remote testdata/spec
	@echo "Spec updated. Review changes with: git diff testdata/spec"

clean:
	go clean ./...
