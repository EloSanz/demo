# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Primero copiamos dependencias para cachear capas
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Stage 2: Final minimal image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

ENTRYPOINT ["./api"]
