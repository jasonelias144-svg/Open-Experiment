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
// GraphTelemetry packages non-linear topological states for display viewportstype GraphTelemetry struct {
	Timestamp      string    `json:"timestamp"`
	LifecycleTick  int64     `json:"lifecycleTick"`
	NodesPotential []float64 `json:"nodesPotential"` // Metabolic activation states
	EdgeMatrix     [][]int   `json:"edgeMatrix"`     // Dynamic connection arrays
	NetworkEntropy float64   `json:"networkEntropy"` // Structural phase variance
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
type MeshNode struct {
	ID        int
	Potential float64
	Neighbors []int
}
type AdaptiveGraph struct {
	Nodes []*MeshNode
	Tick  int64
	Mu    sync.RWMutex
}
func NewAdaptiveGraph(size int) *AdaptiveGraph {
	g := &AdaptiveGraph{
		Nodes: make([]*MeshNode, size),
		Tick:  0,
	}
	for i := 0; i < size; i++ {
		g.Nodes[i] = &MeshNode{
			ID:        i,
			Potential: rand.Float64(),
			Neighbors: []int{},
		}
	}
	// Inject baseline cluster links
	for i := 0; i < size; i++ {
		for j := i + 1; j < size; j++ {
			if rand.Float64() < 0.15 {
				g.Nodes[i].Neighbors = append(g.Nodes[i].Neighbors, j)
				g.Nodes[j].Neighbors = append(g.Nodes[j].Neighbors, i)
			}
		}
	}
	return g
}
func (g *AdaptiveGraph) AdvanceMetabolism() float64 {
	g.Mu.Lock()
	defer g.Mu.Unlock()

	g.Tick++
	size := len(g.Nodes)
	totalEdges := 0

	// Step 1: Compute non-linear activation state shifts
	for _, node := range g.Nodes {
		sumNeighbors := 0.0
		for _, nIdx := range node.Neighbors {
			sumNeighbors += g.Nodes[nIdx].Potential
		}
		
		if len(node.Neighbors) > 0 {
			meanPotential := sumNeighbors / float64(len(node.Neighbors))
			// Sigmoidal transition mapping on localized network states
			node.Potential = 1.0 / (1.0 + math.Exp(-(meanPotential - 0.5) * 4.0))
		}
		
		// Stochastic decay / spontaneous generation factor
		if rand.Float64() < 0.05 {
			node.Potential = rand.Float64()
		}
	}

	// Step 2: Rewrite connection parameters based on phase alignment
	for i := 0; i < size; i++ {
		for j := i + 1; j < size; j++ {
			phaseAlign := math.Abs(g.Nodes[i].Potential - g.Nodes[j].Potential)
			isLinked := contains(g.Nodes[i].Neighbors, j)

			if isLinked && phaseAlign > 0.45 && rand.Float64() < 0.2 {
				// Explode / Prune connection due to structural asynchronous drift
				g.Nodes[i].Neighbors = remove(g.Nodes[i].Neighbors, j)
				g.Nodes[j].Neighbors = remove(g.Nodes[j].Neighbors, i)
			} else if !isLinked && phaseAlign < 0.12 && rand.Float64() < 0.1 {
				// Spawn structural link path due to harmonic convergence
				g.Nodes[i].Neighbors = append(g.Nodes[i].Neighbors, j)
				g.Nodes[j].Neighbors = append(g.Nodes[j].Neighbors, i)
			}
		}
		totalEdges += len(g.Nodes[i].Neighbors)
	}

	entropy := float64(totalEdges) / float64(size*size)
	return entropy
}
func contains(arr []int, val int) bool {
	for _, v := range arr { if v == val { return true } }
	return false
}
func remove(arr []int, val int) []int {
	out := []int{}
	for _, v := range arr { if v != val { out = append(out, v) } }
	return out
}
func main() {
	graph := NewAdaptiveGraph(40) // 40-Node Evolving Micro-Mesh

	http.HandleFunc("/graph-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil { return }
		defer conn.Close()
		fmt.Println("[GRAPH_ENGINE] Telemetry receiver attached to topology core.")

		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			entropy := graph.AdvanceMetabolism()

			graph.Mu.RLock()
			potentials := make([]float64, len(graph.Nodes))
			edges := make([][]int, len(graph.Nodes))
			for i, node := range graph.Nodes {
				potentials[i] = node.Potential
				edges[i] = make([]int, len(node.Neighbors))
				copy(edges[i], node.Neighbors)
			}
			tick := graph.Tick
			graph.Mu.RUnlock()

			packet := GraphTelemetry{
				Timestamp:      time.Now().Format(time.RFC3339),
				LifecycleTick:  tick,
				NodesPotential: potentials,
				EdgeMatrix:     edges,
				NetworkEntropy: entropy,
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil { return }
		}
	})

	fmt.Println("[INIT] Graph Cellular Automata Core broadcasting on ws://localhost:8080/graph-stream")
	_ = http.ListenAndServe(":8080", nil)
}
