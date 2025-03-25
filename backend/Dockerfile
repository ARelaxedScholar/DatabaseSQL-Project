FROM golang:1.24.1 as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

ENV GIN_MODE=release
EXPOSE 8080
CMD ["./main"]
