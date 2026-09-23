# Stage 1: Native binary compilation inside lightweight Linux containerFROM golang:1.21-alpine AS exchange-compiler
WORKDIR /usr/src/market_core
# Copy tracking manifest schemaCOPY go.mod ./RUN go mod download || touch go.mod
# Inject structural processing calculations architecture filesCOPY main.go .
# Execute static binary compiler rules targeted for isolated executionRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o high_frequency_market_engine .
# Stage 2: Final runtime image deployment container constructionFROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layersCOPY --from=exchange-compiler /usr/src/market_core/high_frequency_market_engine .
# Open communications array port for streaming telemetry downstreamEXPOSE 8080
# Launch matching order book clockwork loop on bootCMD ["./high_frequency_market_engine"]
