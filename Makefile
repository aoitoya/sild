BINARY_NAME=sild
INSTALL_DIR=/usr/local/bin

.PHONY: build install uninstall clean

# Build the binary
build:
	go build -o $(BINARY_NAME) ./cmd/app/

# Install the binary to /usr/local/bin
install: build
	sudo install -m 755 $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

# Remove the installed binary
uninstall:
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)