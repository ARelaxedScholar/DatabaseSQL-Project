FROM golang:1.24.1 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Declare the build argument with a default value (if desired)
ARG _ENV_FILE="default_env_value"
RUN echo "$_ENV_FILE" > .env

ENV GIN_MODE=release
EXPOSE 8080
CMD ["./main"]
