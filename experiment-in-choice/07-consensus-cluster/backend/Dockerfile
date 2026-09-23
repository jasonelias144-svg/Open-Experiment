# Stage 1: Native binary compilation inside lightweight Linux containerFROM golang:1.21-alpine AS consensus-compiler
WORKDIR /usr/src/consensus_core
# Copy structural manifest tracking fileCOPY go.mod ./RUN go mod download || touch go.mod
# Inject structural network processing calculations architecture filesCOPY main.go .
# Execute static binary compiler rules targeted for isolated executionRUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consensus_cluster_engine .
# Stage 2: Final runtime image deployment container constructionFROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Pull compiled machine artifacts back out of compiler container layersCOPY --from=consensus-compiler /usr/src/consensus_core/consensus_cluster_engine .
# Open communications array port for streaming consensus data frames downstreamEXPOSE 8080
# Launch decentralized replication clockwork loop on bootCMD ["./consensus_cluster_engine"]
