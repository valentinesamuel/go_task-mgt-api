# Define Go-related variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GORUN = $(GOCMD) run
GOTIDY = $(GOCMD) mod tidy

# Paths
MAIN_PATH = ./main.go
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

.PHONY: start build test clean lint up logs app-logs elk-logs down restart-logging check-indices

up:
	@echo "Starting all services..."
	docker-compose up -d

logs:
	@echo "Showing logs for all services..."
	docker-compose logs -f

app-logs:
	@echo "Showing application logs..."
	tail -f logs/*.log

elk-logs:
	@echo "Showing Fluent Bit logs..."
	docker-compose logs -f fluentbit

down:
	@echo "Stopping all services..."
	docker-compose down

clean:
	@echo "Cleaning up unused resources across all services..."
	docker-compose down -v --remove-orphans

restart-logging:
	@echo "Restarting Fluent Bit service..."
	docker-compose restart fluentbit

check-indices:
	@echo "Checking Elasticsearch indices..."
	curl -X GET "http://localhost:9200/_cat/indices?v"
