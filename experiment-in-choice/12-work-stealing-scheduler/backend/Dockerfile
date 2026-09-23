# Stage 1: Fast machine compilation box using minimal Go environmentFROM golang:1.21-alpine AS scheduler-compilerWORKDIR /usr/src/scheduler_matrix
# Cache tracking descriptors for swift image rebuild passesCOPY go.mod go.sum* ./RUN go mod download || touch go.mod
# Inject the core matrix generation logic filesCOPY scheduling_engine.go .
# Execute compiler constraints to output an unlinked static binary assetRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o work_stealing_engine .
# Stage 2: Hardened runtime layout with zero external packagesFROM alpine:latest  RUN apk --no-cache add ca-certificatesWORKDIR /root/
# Reclaim the compiled engine block from the compilation layerCOPY --from=scheduler-compiler /usr/src/scheduler_matrix/work_stealing_engine .
# Open specific TCP port allocations for low-latency WebSocket streamingEXPOSE 8080
# Trigger the lock-free processing loop automatically upon box initializationCMD ["./work_stealing_engine"]
