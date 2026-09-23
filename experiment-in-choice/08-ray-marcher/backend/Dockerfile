# Stage 1: Native binary compilation inside lightweight Linux containerFROM golang:1.21-alpine AS graphics-compiler
WORKDIR /usr/src/ray_core
# Copy structural manifest tracking fileCOPY go.mod ./RUN go mod download || touch go.mod
# Inject structural graphics processing calculations architecture filesCOPY main.go .
# Execute static binary compiler rules targeted for isolated executionRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ray_marching_engine .
# Stage 2: Final runtime image deployment container constructionFROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layersCOPY --from=graphics-compiler /usr/src/ray_core/ray_marching_engine .
# Open communications array port for streaming frame buffers downstreamEXPOSE 8080
# Launch graphics rasterization clockwork loop on bootCMD ["./ray_marching_engine"]
