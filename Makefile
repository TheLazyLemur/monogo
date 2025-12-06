.PHONY: release patch minor major build test clean

# Get current version from git tags
CURRENT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
VERSION_PARTS := $(subst ., ,$(subst v,,$(CURRENT_VERSION)))
MAJOR := $(word 1,$(VERSION_PARTS))
MINOR := $(word 2,$(VERSION_PARTS))
PATCH := $(word 3,$(VERSION_PARTS))

# Calculate next versions
NEXT_PATCH := v$(MAJOR).$(MINOR).$(shell echo $$(($(PATCH)+1)))
NEXT_MINOR := v$(MAJOR).$(shell echo $$(($(MINOR)+1))).0
NEXT_MAJOR := v$(shell echo $$(($(MAJOR)+1))).0.0

help:
	@echo "MonoGo Release Management"
	@echo ""
	@echo "Current version: $(CURRENT_VERSION)"
	@echo ""
	@echo "Commands:"
	@echo "  make patch    - Release $(NEXT_PATCH) (bug fixes)"
	@echo "  make minor    - Release $(NEXT_MINOR) (new features)"
	@echo "  make major    - Release $(NEXT_MAJOR) (breaking changes)"
	@echo "  make build    - Build all packages"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo ""
	@echo "Example: make patch"

build:
	go build ./...

test:
	go test ./...

clean:
	rm -f example_rts/example_rts
	go clean ./...

patch: build test
	@echo "Releasing $(NEXT_PATCH)..."
	git tag $(NEXT_PATCH)
	git push origin $(NEXT_PATCH)
	@echo "Released $(NEXT_PATCH)"

minor: build test
	@echo "Releasing $(NEXT_MINOR)..."
	git tag $(NEXT_MINOR)
	git push origin $(NEXT_MINOR)
	@echo "Released $(NEXT_MINOR)"

major: build test
	@echo "Releasing $(NEXT_MAJOR)..."
	@echo "WARNING: Major version bump. Make sure breaking changes are intentional."
	@read -p "Continue? [y/N] " confirm && [ "$$confirm" = "y" ]
	git tag $(NEXT_MAJOR)
	git push origin $(NEXT_MAJOR)
	@echo "Released $(NEXT_MAJOR)"

# Release specific version: make release V=v1.2.3
release: build test
ifndef V
	$(error Usage: make release V=v1.2.3)
endif
	@echo "Releasing $(V)..."
	git tag $(V)
	git push origin $(V)
	@echo "Released $(V)"
