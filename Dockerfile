# Stage 1: Builder
FROM golang:1.22-alpine AS builder

# Add necessary build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main .

# Stage 2: Final executable
FROM alpine:latest

# Add CA certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates

# Create a non-root user
RUN adduser -D appuser

WORKDIR /app

# Copy the executable from the builder stage
COPY --from=builder /app/main .

# Copy .env file
COPY .env ./.env

# Set ownership of the application binary
RUN chown appuser:appuser /app/main

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 3000

# Set environment variables (instead of .env file)
ENV GIN_MODE=release

CMD ["./main"]