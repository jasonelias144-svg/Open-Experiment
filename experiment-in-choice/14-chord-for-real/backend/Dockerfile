FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o chord .

FROM alpine:latest
COPY --from=builder /src/chord /chord
EXPOSE 8080
CMD ["/chord"]
