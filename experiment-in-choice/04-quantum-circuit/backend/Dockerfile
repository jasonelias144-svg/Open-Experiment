# Stage 1: Native binary compilation inside lightweight Linux container
FROM golang:1.21-alpine AS quantum-compiler
WORKDIR /usr/src/quantum_core
# Copy structural manifest tracking file
COPY go.mod go.sum ./
RUN go mod download || touch go.mod
# Inject structural matrix processing logic files
COPY main.go .
# Execute static binary compiler rules targeted for isolated execution
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quantum_simulator_engine .
# Stage 2: Final runtime image deployment container construction
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layers
COPY --from=quantum-compiler /usr/src/quantum_core/quantum_simulator_engine .
# Open communications array port for streaming quantum data frames downstream
EXPOSE 8080
# Launch coherent state computation clockwork loop on boot
CMD ["./quantum_simulator_engine"]
