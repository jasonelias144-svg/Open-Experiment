package main
import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"://github.com"
)
const RingModulo = 1000 // Simplified virtual token ring bounds for presentation mapping
// ChordTelemetry packages decentralized routing structures for web display viewportstype ChordTelemetry struct {
	Timestamp      string           `json:"timestamp"`
	LifecycleTick  int64            `json:"lifecycleTick"`
	NodePositions  []int            `json:"nodePositions"`  // Active node positions on 0-999 ring
	KeyAssignments map[int][]string `json:"keyAssignments"` // NodePosition -> Key Hashes allocated
	RoutingSteps   int              `json:"routingSteps"`   // Hops spent on last content lookup
	NetworkEntropy float64          `json:"networkEntropy"` // Variance profile of load distribution
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
type ChordNode struct {
	HashID      int
	AssignedKeys []string
}
type ChordRing struct {
	Nodes map[int]*ChordNode
	Keys  map[string]int // Raw string key -> its absolute integer hash location
	Tick  int64
	Mu    sync.RWMutex
}
func NewChordRing() *ChordRing {
	return &ChordRing{
		Nodes: make(map[int]*ChordNode),
		Keys:  make(map[string]int),
		Tick:  0,
	}
}
func hashString(key string) int {
	h := sha1.New()
	h.Write([]byte(key))
	bs := h.Sum(nil)
	val := binary.BigEndian.Uint32(bs[0:4])
	return int(val % uint32(RingModulo))
}
func (r *ChordRing) AddNode(nodeName string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	hID := hashString(nodeName)
	if _, exists := r.Nodes[hID]; !exists {
		r.Nodes[hID] = &ChordNode{HashID: hID, AssignedKeys: []string{}}
		fmt.Printf("[DHT_RING] Mounted storage instance node: %s at token coordinate: %d\n", nodeName, hID)
	}
}
func (r *ChordRing) RemoveNodeRandom() {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	if len(r.Nodes) <= 2 {
		return // Retain a critical quorum floor for mapping integrity
	}
	
	// Select random key position to prune
	keys := make([]int, 0, len(r.Nodes))
	for k := range r.Nodes {
		keys = append(keys, k)
	}
	target := keys[rand.Intn(len(keys))]
	delete(r.Nodes, target)
	fmt.Printf("[!!! CHURN] Node instance coordinate %d dropped completely offline.\n", target)
}
func (r *ChordRing) ReallocateTopology() (int, float64) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	r.Tick++

	// Extract and sort active token positions
	var sortedNodes []int
	for k := range r.Nodes {
		r.Nodes[k].AssignedKeys = []string{} // Reset buffer maps for re-evaluation pass
		sortedNodes = append(sortedNodes, k)
	}
	sort.Ints(sortedNodes)

	// Inject periodic fresh asset blocks into the data partition ring
	if r.Tick%5 == 0 {
		mockAssetID := fmt.Sprintf("content_hash_0x%X", rand.Int64())
		r.Keys[mockAssetID] = hashString(mockAssetID)
	}

	// Consistent Hashing mapping: Route keys to the closest successor node position
	for keyStr, keyHash := range r.Keys {
		successorIdx := 0
		found := false
		for i, nodePos := range sortedNodes {
			if nodePos >= keyHash {
				successorIdx = i
				found = true
				break
			}
		}
		if !found {
			successorIdx = 0 // Wrap around boundaries back to the first ring successor
		}
		targetNodePos := sortedNodes[successorIdx]
		r.Nodes[targetNodePos].AssignedKeys = append(r.Nodes[targetNodePos].AssignedKeys, keyStr)
	}

	// Calculate a mock hop tracking index approximating look-ahead tables logs
	hopsSpent := 1 + rand.Intn(int(math.Log2(float64(len(sortedNodes)+1)))+1)

	// Calculate network distribution entropy balance
	var averageLoad = float64(len(r.Keys)) / float64(len(sortedNodes))
	varianceSum := 0.0
	for _, n := range r.Nodes {
		diff := float64(len(n.AssignedKeys)) - averageLoad
		varianceSum += diff * diff
	}
	entropyMetric := math.Sqrt(varianceSum / float64(len(sortedNodes)))

	return hopsSpent, entropyMetric
}
func main() {
	ring := NewChordRing()
	
	// Seed static network topology ring nodes
	ring.AddNode("storage_cell_alpha")
	ring.AddNode("storage_cell_beta")
	ring.AddNode("storage_cell_gamma")
	ring.AddNode("storage_cell_delta")
	ring.AddNode("storage_cell_epsilon")

	http.HandleFunc("/chord-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[CHORD] Observer socket handshaked with consistent hashing cluster.")

		ticker := time.NewTicker(200 * time.Millisecond) // Evaluation clock frequency loop
		defer ticker.Stop()

		for range ticker.C {
			// Simulate live cluster network churn dynamics
			if rand.Float64() < 0.04 {
				ring.RemoveNodeRandom()
			}
			if rand.Float64() < 0.05 {
				ring.AddNode(fmt.Sprintf("dynamic_cell_%d", rand.Intn(1000)))
			}

			hops, entropy := ring.ReallocateTopology()

			ring.Mu.RLock()
			nodePositions := make([]int, 0, len(ring.Nodes))
			keyAssignments := make(map[int][]string)
			
			for pos, node := range ring.Nodes {
				nodePositions = append(nodePositions, pos)
				keyAssignments[pos] = make([]string, len(node.AssignedKeys))
				copy(keyAssignments[pos], node.AssignedKeys)
			}
			tick := ring.Tick
			ring.Mu.RUnlock()

			packet := ChordTelemetry{
				Timestamp:      time.Now().Format(time.RFC3339),
				LifecycleTick:  tick,
				NodePositions:  nodePositions,
				KeyAssignments: keyAssignments,
				RoutingSteps:   hops,
				NetworkEntropy: entropy,
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Consistent Hashing Chord Ring active on ws://localhost:8080/chord-stream")
	_ = http.ListenAndServe(":8080", nil)
}
