# Start from the official Golang image for build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd/main.go

# Use a minimal base image for running
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth .
EXPOSE 8080
CMD ["./auth"] 