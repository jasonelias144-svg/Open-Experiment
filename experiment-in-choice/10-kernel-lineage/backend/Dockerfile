# Stage 1: Build the Go SQLite binary with full GCC compilation bindings
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
# Add and initialize necessary sqlite dependencies module references
RUN go get github.com/mattn/go-sqlite3 && go mod tidy || touch go.mod
COPY main.go .
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main .
# Stage 2: Deploy execution layer
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
# Expose the communication array ports
EXPOSE 8080
CMD ["./main"]
