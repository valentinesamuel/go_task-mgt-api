# Define Go-related variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GORUN = $(GOCMD) run
GOTIDY = $(GOCMD) mod tidy

# Paths
MAIN_PATH = ./cmd/api/main.go
BIN_PATH = ./bin/myapp

# Default target to build and start the application
start:
	$(GORUN) $(MAIN_PATH)

# Build the application into a binary
build:
	$(GOBUILD) -o $(BIN_PATH) $(MAIN_PATH)

# Run tests
test:
	$(GOTEST) ./...

# Clean up binaries and other generated files
clean:
	$(GOCLEAN)
	rm -rf $(BIN_PATH)

# Tidy up Go modules
tidy:
	$(GOTIDY)

# Run the linter (if using a linter like golangci-lint)
lint:
	golangci-lint run

.PHONY: start build test clean lint
