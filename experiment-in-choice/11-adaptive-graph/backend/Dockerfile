# Stage 1: Fast machine compilation box using minimal Go environment
FROM golang:1.21-alpine AS engine-builder
WORKDIR /usr/src/graph_matrix
# Cache tracking descriptors for swift image rebuild passes
COPY go.mod go.sum* ./
RUN go mod download || touch go.mod
# Inject the core matrix generation logic files
COPY dynamic_graph_engine.go .
# Execute compiler constraints to output an unlinked static binary asset
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graph_automata_engine .
# Stage 2: Hardened runtime layout with zero external packages
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Reclaim the compiled engine block from the compilation layer
COPY --from=engine-builder /usr/src/graph_matrix/graph_automata_engine .
# Open specific TCP port allocations for low-latency WebSocket streaming
EXPOSE 8080
# Trigger the metabolic generation loop automatically upon box initialization
CMD ["./graph_automata_engine"]
