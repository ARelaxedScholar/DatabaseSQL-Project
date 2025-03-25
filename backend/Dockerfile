# Dockerfile
# Use Go base image
FROM golang:1.24.1 as builder

WORKDIR /app

# Copy all files to the container
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./main.go

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Set environment variable
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./main"]
