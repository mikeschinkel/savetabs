# Variables
DAEMON_BIN_DIR = ./daemon/bin
DAEMON_CMD_DIR = ./daemon/cmd
DAEMON_SRC = savetabsd/main.go
DAEMON_BIN = ./daemon/bin/savetabsd
DAEMON_DIR = ./daemon

# Run the daemon
run-daemon: build-daemon
	cd $(DAEMON_DIR) && ./bin/savetabsd

# Build the daemon binary
build-daemon:
	mkdir -p $(DAEMON_BIN_DIR)
	cd $(DAEMON_CMD_DIR) \
	  	&& go mod tidy \
		&& go build -o $(DAEMON_BIN) $(DAEMON_SRC)

# Clean up the build artifacts
clean-daemon:
	rm -rf $(DAEMON_DIR)/bin

.PHONY: build-daemon run-daemon clean-daemon
