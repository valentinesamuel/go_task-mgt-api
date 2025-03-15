# Logging Architecture

This document explains how logs flow from the application to Elasticsearch using Fluent Bit.

## Log Flow

1. **Application Logging**
   - The application uses Logrus with ECS formatter to generate structured logs
   - Logs are written to both stdout and the `./logs/out.log` file
   - The ECS formatter ensures logs are in Elasticsearch Common Schema format

2. **Fluent Bit Collection**
   - Fluent Bit monitors the `./logs/*.log` files for changes
   - It uses the `tail` input plugin to read new log entries
   - Logs are parsed using the JSON parser with ECS format support

3. **Log Processing**
   - Fluent Bit processes the logs according to the configuration
   - The `@timestamp` field from ECS is used for the log timestamp
   - UTF-8 encoded messages are properly decoded

4. **Elasticsearch Output**
   - Logs are sent to Elasticsearch running at `elasticsearch:9200`
   - They are indexed with the pattern `go-task-api-%Y.%m.%d` (e.g., `go-task-api-2025.03.15`)
   - No additional processing is done at the Elasticsearch level

5. **Visualization in Kibana**
   - Logs can be viewed and analyzed in Kibana at `http://localhost:5601`
   - You can create dashboards and visualizations based on the log data

## Useful Commands

Use these Makefile commands to manage the logging system:

```bash
# Start all services including Fluent Bit and Elasticsearch
make up

# View Fluent Bit logs to troubleshoot log shipping
make elk-logs

# View application logs directly
make app-logs

# Check Elasticsearch indices to verify log delivery
make check-indices

# Restart just the logging service if needed
make restart-logging
```

## Configuration Files

- **fluent-bit.conf**: Main configuration for Fluent Bit
- **parsers.conf**: Defines how to parse the JSON logs with ECS format

## Troubleshooting

If logs are not appearing in Elasticsearch:

1. Check if application logs are being generated: `make app-logs`
2. Verify Fluent Bit is running: `docker-compose ps`
3. Check Fluent Bit logs for errors: `make elk-logs`
4. Ensure Elasticsearch is accessible: `curl http://localhost:9200`
5. Check for existing indices: `make check-indices`
