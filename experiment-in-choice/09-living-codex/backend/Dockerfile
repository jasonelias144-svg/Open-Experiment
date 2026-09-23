FROM golang:1.21-alpine AS codex-compiler
WORKDIR /usr/src/quantum_core
COPY go.mod go.sum ./
RUN go mod download || touch go.mod
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quantum_simulator_engine .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=codex-compiler /usr/src/quantum_core/quantum_simulator_engine .
EXPOSE 8080
CMD ["./quantum_simulator_engine"]
