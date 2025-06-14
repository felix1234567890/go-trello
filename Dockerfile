# Start from the official Golang image as a build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o go-trello main.go

# Start a minimal image for running
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/go-trello .
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./go-trello"]
