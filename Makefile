.PHONY: build clean install uninstall

# Build variables
BINARY_NAME=vocl
BUILD_DIR=build
MAIN_PATH=main.go

# You can override these when calling make, e.g.:
# make build GOOS=linux GOARCH=amd64
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED ?= 1

# Default target
all: build

# Build the application
build: clean
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Build complete. Binary is in $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete"

# Uninstall the binary from /usr/local/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Help target
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install the binary to /usr/local/bin"
	@echo "  uninstall  - Remove the binary from /usr/local/bin"
	@echo "  help       - Show this help message" 