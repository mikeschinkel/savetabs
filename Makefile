# Variables
ROOT_DIR = $(shell pwd)
DAEMON_DIR = daemon
EXTENSION_DIR = extension
DAEMON_BIN_DIR = $(DAEMON_DIR)/bin
DAEMON_CMD_DIR = $(DAEMON_DIR)/cmd/savetabsd
DAEMON_SRC = $(DAEMON_CMD_DIR)/main.go
DAEMON_BIN = $(DAEMON_BIN_DIR)/savetabsd

# Run the daemon
run-daemon: build-daemon build-extension
	cd $(DAEMON_DIR) && ./bin/savetabsd

# Build the daemon binary
build-daemon: clean-daemon
	mkdir -p $(DAEMON_BIN_DIR)
	cd "$(ROOT_DIR)/$(DAEMON_DIR)" \
	  	&& go mod tidy \
		&& go build -o "${ROOT_DIR}/$(DAEMON_BIN)" "${ROOT_DIR}/$(DAEMON_SRC)"

# Build the daemon binary
build-extension:
	cd "$(ROOT_DIR)/$(EXTENSION_DIR)" \
	  	&& make build

# Clean up the build artifacts
clean-daemon:
	rm -rf "$(ROOT_DIR)/$(DAEMON_BIN_DIR)"

.PHONY: build-daemon build-extension run-daemon clean-daemon
