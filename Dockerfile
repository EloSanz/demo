# Stage 1: Build the application
FROM golang:1.23.4-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a static binary for the final minimal image
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Stage 2: Run the application
FROM alpine:latest

# Install CA certificates for external HTTPS calls (JSONPlaceholder, Rick and Morty)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/api .

# Expose the port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./api"]
