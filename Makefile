# Variables
ROOT_DIR = $(shell pwd)
DAEMON_DIR = daemon
DAEMON_BIN_DIR = $(DAEMON_DIR)/bin
DAEMON_CMD_DIR = $(DAEMON_DIR)/cmd/savetabsd
DAEMON_SRC = $(DAEMON_CMD_DIR)/main.go
DAEMON_BIN = $(DAEMON_BIN_DIR)/savetabsd

# Run the daemon
run-daemon: build-daemon
	cd $(DAEMON_DIR) && ./bin/savetabsd

# Build the daemon binary
build-daemon: clean-daemon
	mkdir -p $(DAEMON_BIN_DIR)
	cd "$(ROOT_DIR)/$(DAEMON_DIR)" \
	  	&& go mod tidy \
		&& go build -o "${ROOT_DIR}/$(DAEMON_BIN)" "${ROOT_DIR}/$(DAEMON_SRC)"

# Clean up the build artifacts
clean-daemon:
	rm -rf "$(ROOT_DIR)/$(DAEMON_BIN_DIR)"

.PHONY: build-daemon run-daemon clean-daemon
