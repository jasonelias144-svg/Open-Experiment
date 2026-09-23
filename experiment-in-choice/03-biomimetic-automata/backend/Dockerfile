# Stage 1: Native binary compilation inside lightweight Linux containerFROM golang:1.21-alpine AS engine-compiler
WORKDIR /usr/src/bio_matrix
# Copy tracking manifest schemas for local modules initializationCOPY go.mod ./RUN go mod download || touch go.mod
# Inject structural processing calculations architecture filesCOPY main.go .
# Execute static binary compiler rules targeted for isolated executionRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bio_cellular_engine .
# Stage 2: Final runtime image deployment container constructionFROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layersCOPY --from=engine-compiler /usr/src/bio_matrix/bio_cellular_engine .
# Open communications array port for streaming data frames downstreamEXPOSE 8080
# Spin up metabolic simulation clockwork loop automatically on bootCMD ["./bio_cellular_engine"]
