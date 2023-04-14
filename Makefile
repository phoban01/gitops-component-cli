# Set the name of your binary
BINARY_NAME := component

# Set the directory where your source code is located
SRC_DIR := ./cmd/component

# Set the Go compiler command
GO := go

# Set the build flags
BUILD_FLAGS := -ldflags="-s -w"

# Set the installation directory
INSTALL_DIR := /usr/local/bin

# Set the build directory
BUILD_DIR := ./bin

# Define the default target
.PHONY: build
build: $(BUILD_DIR)/$(BINARY_NAME)

# Define the build target
$(BUILD_DIR)/$(BINARY_NAME):
	@$(GO) build $(BUILD_FLAGS) -o $@ $(SRC_DIR)

# Define the install target
.PHONY: install
install: $(BUILD_DIR)/$(BINARY_NAME)
	@cp $< $(INSTALL_DIR)

# Define the clean target
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)
