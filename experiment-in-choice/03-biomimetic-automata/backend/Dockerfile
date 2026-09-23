# Stage 1: Native binary compilation inside lightweight Linux container
FROM golang:1.21-alpine AS engine-compiler
WORKDIR /usr/src/bio_matrix
# Copy tracking manifest schemas for local modules initialization
COPY go.mod go.sum ./
RUN go mod download || touch go.mod
# Inject structural processing calculations architecture files
COPY main.go .
# Execute static binary compiler rules targeted for isolated execution
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bio_cellular_engine .
# Stage 2: Final runtime image deployment container construction
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layers
COPY --from=engine-compiler /usr/src/bio_matrix/bio_cellular_engine .
# Open communications array port for streaming data frames downstream
EXPOSE 8080
# Spin up metabolic simulation clockwork loop automatically on boot
CMD ["./bio_cellular_engine"]
