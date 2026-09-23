# Stage 1: Native binary compilation inside lightweight Linux container
FROM golang:1.21-alpine AS graphics-compiler
WORKDIR /usr/src/ray_core
# Copy structural manifest tracking file
COPY go.mod go.sum ./
RUN go mod download || touch go.mod
# Inject structural graphics processing calculations architecture files
COPY main.go .
# Execute static binary compiler rules targeted for isolated execution
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ray_marching_engine .
# Stage 2: Final runtime image deployment container construction
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layers
COPY --from=graphics-compiler /usr/src/ray_core/ray_marching_engine .
# Open communications array port for streaming frame buffers downstream
EXPOSE 8080
# Launch graphics rasterization clockwork loop on boot
CMD ["./ray_marching_engine"]
