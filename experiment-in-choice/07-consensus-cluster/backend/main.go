package main
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"://github.com"
)
// ConsensusTelemetry encodes cluster sync logs for visualization viewportstype ConsensusTelemetry struct {
	Timestamp     string   `json:"timestamp"`
	CurrentTerm   int64    `json:"currentTerm"`
	LatestBlock   int64    `json:"latestBlock"`
	NodeStates    []string `json:"nodeStates"`    // States: "LEADER" | "FOLLOWER" | "OFFLINE"
	BlockHeights  []int64  `json:"blockHeights"`  // Replicated ledger height per node
	NetworkHealth float64  `json:"networkHealth"` // Active node consensus consensus ratio
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
type ValidatorNode struct {
	ID          int
	State       string // "LEADER", "FOLLOWER", "OFFLINE"
	BlockHeight int64
}
type ConsensusCluster struct {
	Nodes       []*ValidatorNode
	Term        int64
	GlobalBlock int64
	Mu          sync.Mutex
}
func NewCluster() *ConsensusCluster {
	c := &ConsensusCluster{
		Nodes: make([]*ValidatorNode, 5), // 5-Node consensus topology
		Term:  1,
	}
	for i := 0; i < 5; i++ {
		c.Nodes[i] = &ValidatorNode{ID: i, State: "FOLLOWER", BlockHeight: 0}
	}
	c.Nodes[0].State = "LEADER" // Node 0 starts as primary proposer
	return c
}
func (c *ConsensusCluster) ProcessNetworkTick() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	activeQuorum := 0
	leaderAlive := false

	// Inject randomized network dropping faults to simulate Byzantine/crash failures
	for _, node := range c.Nodes {
		if rand.Float64() < 0.12 {
			node.State = "OFFLINE"
		} else if node.State == "OFFLINE" {
			node.State = "FOLLOWER"
		}
		
		if node.State == "LEADER" {
			leaderAlive = true
		}
		if node.State != "OFFLINE" {
			activeQuorum++
		}
	}

	// Trigger emergency re-election loop if leader dropped offline
	if !leaderAlive && activeQuorum >= 3 {
		c.Term++
		for _, node := range c.Nodes {
			if node.State != "OFFLINE" {
				node.State = "LEADER" // Simple fallback: first online node takes term lead
				fmt.Printf("[CONSENSUS] Leader crashed. Term %d election complete.\n", c.Term)
				break
			}
		}
	}

	// If quorum exists, leader appends a block and propagates state entries
	if activeQuorum >= 3 {
		c.GlobalBlock++
		for _, node := range c.Nodes {
			if node.State != "OFFLINE" {
				node.BlockHeight = c.GlobalBlock
			}
		}
	}
}
func main() {
	cluster := NewCluster()

	http.HandleFunc("/consensus-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[CONSENSUS] Ledger diagnostic monitor linked to cluster fabric.")

		ticker := time.NewTicker(350 * time.Millisecond) // Consensus tick step interval
		defer ticker.Stop()

		for range ticker.C {
			cluster.ProcessNetworkTick()

			states := make([]string, 5)
			heights := make([]int64, 5)
			onlineCount := 0.0

			cluster.Mu.Lock()
			for i, n := range cluster.Nodes {
				states[i] = n.State
				heights[i] = n.BlockHeight
				if n.State != "OFFLINE" {
					onlineCount++
				}
			}
			
			packet := ConsensusTelemetry{
				Timestamp:     time.Now().Format(time.RFC3339),
				CurrentTerm:   cluster.Term,
				LatestBlock:   cluster.GlobalBlock,
				NodeStates:    states,
				BlockHeights:  heights,
				NetworkHealth: (onlineCount / 5.0) * 100.0,
			}
			cluster.Mu.Unlock()

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Consensus Simulation cluster spinning on ws://localhost:8080/consensus-stream")
	_ = http.ListenAndServe(":8080", nil)
}
