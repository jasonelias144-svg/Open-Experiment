# Stage 1: Fast machine compilation box using minimal Go environmentFROM golang:1.21-alpine AS ring-compilerWORKDIR /usr/src/dht_matrix
# Cache tracking descriptors for swift image rebuild passesCOPY go.mod go.sum* ./RUN go mod download || touch go.mod
# Inject the core matrix generation logic filesCOPY dht_chord_engine.go .
# Execute compiler constraints to output an unlinked static binary assetRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dht_chord_engine .
# Stage 2: Hardened runtime layout with zero external packagesFROM alpine:latest  RUN apk --no-cache add ca-certificatesWORKDIR /root/
# Reclaim the compiled engine block from the compilation layerCOPY --from=ring-compiler /usr/src/dht_matrix/dht_chord_engine .
# Open specific TCP port allocations for low-latency WebSocket streamingEXPOSE 8080
# Trigger the decentralized storage loop automatically upon box initializationCMD ["./dht_chord_engine"]
