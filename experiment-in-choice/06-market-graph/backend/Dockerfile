# Stage 1: Native binary compilation inside lightweight Linux container
FROM golang:1.21-alpine AS exchange-compiler
WORKDIR /usr/src/market_core
# Copy tracking manifest schema
COPY go.mod go.sum ./
RUN go mod download || touch go.mod
# Inject structural processing calculations architecture files
COPY main.go .
# Execute static binary compiler rules targeted for isolated execution
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o high_frequency_market_engine .
# Stage 2: Final runtime image deployment container construction
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layers
COPY --from=exchange-compiler /usr/src/market_core/high_frequency_market_engine .
# Open communications array port for streaming telemetry downstream
EXPOSE 8080
# Launch matching order book clockwork loop on boot
CMD ["./high_frequency_market_engine"]
