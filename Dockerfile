FROM golang:1.23.2-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

# Set CGO_ENABLED=1 for SQLite
ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled
RUN go build -o main .

# Production stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    sqlite \
    sqlite-libs

WORKDIR /app

# Copy binary and configs
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env .
COPY --from=builder /app/config ./config

EXPOSE 8080

# Run the binary instead of go run
CMD ["./main"]