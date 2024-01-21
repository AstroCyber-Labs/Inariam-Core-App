# List of all commands
CMDS=$(notdir $(wildcard cmd/*))

# Default Go linker flags
LDFLAGS=-ldflags "-s -w"

# Binary output directory
BIN_DIR=./bin

# Default make target
.DEFAULT_GOAL := help

# Build all commands
build: $(CMDS)

$(CMDS):
	@echo "Building $@..."
	@go build $(LDFLAGS) -o $(BIN_DIR)/$@ ./cmd/$@

# Run tests for a specific package
test-pkg-%:
	@go test ./pkgs/$*

test-cloud-%:
	@echo "Testing ./pkgs/cloud/$*"
	@go test ./pkgs/cloud/$*
# Clean up binaries
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

dev: swagger build
	@docker-compose -f ./scripts/dev.docker-compose.yml up -d --build
	@#./bin/inariam server


swagger:
	@swag init -d core/services/api -g api.go -o core/services/api/docs

# List available make targets
help:
	@echo "Choose a command to run:"
	@echo "  make build (db-migrator,gen-services,inariam)        # Build all binaries under cmd/"
	@echo "  make <cmd_name>       # Build a specific binary under cmd/"
	@echo "  make test-pkg-<pkg>   # Run tests for a specific package under pkg/"
	@echo "  make test-cloud-<pkg>   # Run tests for a specific package under pkgs/cloud/"
	@echo "  make clean            # Remove all built binaries"

.PHONY: build clean help dev swagger
