.PHONY: all install clean

BINARY_NAME=went
BUILD_DIR=./build
INSTALL_DIR=/usr/local/bin
SOURCE_FILES=$(shell find . -name "*.go" -type f)

all: build

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(SOURCE_FILES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) installed to $(INSTALL_DIR)"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "✓ Build artifacts cleaned"

uninstall:
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) uninstalled"

help:
	@echo "Usage:"
	@echo "  make build       Build the project"
	@echo "  make install     Install the binary to $(INSTALL_DIR)"
	@echo "  make clean       Remove build artifacts"
	@echo "  make uninstall   Remove the installed binary"
	@echo "  make help        Show this help message"
