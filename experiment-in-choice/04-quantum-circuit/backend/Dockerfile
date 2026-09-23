# Stage 1: Native binary compilation inside lightweight Linux containerFROM golang:1.21-alpine AS quantum-compiler
WORKDIR /usr/src/quantum_core
# Copy structural manifest tracking fileCOPY go.mod ./RUN go mod download || touch go.mod
# Inject structural matrix processing logic filesCOPY main.go .
# Execute static binary compiler rules targeted for isolated executionRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quantum_simulator_engine .
# Stage 2: Final runtime image deployment container constructionFROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layersCOPY --from=quantum-compiler /usr/src/quantum_core/quantum_simulator_engine .
# Open communications array port for streaming quantum data frames downstreamEXPOSE 8080
# Launch coherent state computation clockwork loop on bootCMD ["./quantum_simulator_engine"]
