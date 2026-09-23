# Stage 1: Fast machine compilation box using minimal Go environment
FROM golang:1.21-alpine AS scheduler-compiler
WORKDIR /usr/src/scheduler_matrix
# Cache tracking descriptors for swift image rebuild passes
COPY go.mod go.sum* ./
RUN go mod download || touch go.mod
# Inject the core matrix generation logic files
COPY scheduling_engine.go .
# Execute compiler constraints to output an unlinked static binary asset
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o work_stealing_engine .
# Stage 2: Hardened runtime layout with zero external packages
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Reclaim the compiled engine block from the compilation layer
COPY --from=scheduler-compiler /usr/src/scheduler_matrix/work_stealing_engine .
# Open specific TCP port allocations for low-latency WebSocket streaming
EXPOSE 8080
# Trigger the lock-free processing loop automatically upon box initialization
CMD ["./work_stealing_engine"]
