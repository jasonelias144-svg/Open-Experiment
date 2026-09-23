package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// NeuralTelemetry encodes synaptic mapping arrays for visualization viewports
type NeuralTelemetry struct {
	Timestamp      string    `json:"timestamp"`
	Generation     int64     `json:"generation"`
	FitnessScore   float64   `json:"fitnessScore"`
	NodeStates     []float64 `json:"nodeStates"`     // Node activation potentials: σ(Σwx)
	SynapseWeights []float64 `json:"synapseWeights"` // Connection multipliers
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Synapse represents a directed learning arc between two nodes
type Synapse struct {
	From   int
	To     int
	Weight float64
}
type NeuralGraph struct {
	Nodes    []float64
	Synapses []Synapse
	Fitness  float64
}

func NewGraph() *NeuralGraph {
	g := &NeuralGraph{
		Nodes: make([]float64, 6), // 2 Input, 2 Hidden, 2 Output architecture
		Synapses: []Synapse{
			{0, 2, rand.Float64()*2 - 1}, {0, 3, rand.Float64()*2 - 1},
			{1, 2, rand.Float64()*2 - 1}, {1, 3, rand.Float64()*2 - 1},
			{2, 4, rand.Float64()*2 - 1}, {2, 5, rand.Float64()*2 - 1},
			{3, 4, rand.Float64()*2 - 1}, {3, 5, rand.Float64()*2 - 1},
		},
	}
	return g
}

// Sigmoid activation clamping function
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
func (g *NeuralGraph) Forward(inputA, inputB float64) {
	g.Nodes[0] = inputA
	g.Nodes[1] = inputB

	// Reset activation tracking potentials for downstream nodes
	for i := 2; i < len(g.Nodes); i++ {
		g.Nodes[i] = 0.0
	}

	// Accumulate inputs scaled across current synaptic weights
	for _, syn := range g.Synapses {
		g.Nodes[syn.To] += g.Nodes[syn.From] * syn.Weight
	}

	// Pass calculated weights through sigmoid activation caps
	for i := 2; i < len(g.Nodes); i++ {
		g.Nodes[i] = sigmoid(g.Nodes[i])
	}

	// Mock fitness calculation evaluation logic: how well outputs target an XOR baseline
	targetXOR := math.Abs(inputA - inputB)
	errorDelta := math.Abs(g.Nodes[4] - targetXOR)
	g.Fitness = 1.0 / (errorDelta + 0.01)
}
func (g *NeuralGraph) PerturbWeights() {
	// Inject gradient-free mutation steps to simulate evolutionary descent optimization paths
	for i := range g.Synapses {
		if rand.Float64() < 0.3 {
			g.Synapses[i].Weight += rand.NormFloat64() * 0.15
		}
	}
}
func main() {
	graph := NewGraph()

	http.HandleFunc("/neural-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[NEURAL] Analytics monitor terminal mounted onto network core.")

		ticker := time.NewTicker(250 * time.Millisecond) // Generation tracking cycle step
		defer ticker.Stop()

		var gen int64 = 0

		for range ticker.C {
			gen++

			// Cycle inputs between raw logic parameters to assess ongoing map efficiency
			inA := 0.0
			inB := 0.0
			if gen%2 == 0 {
				inA = 1.0
			}
			if gen%4 < 2 {
				inB = 1.0
			}

			graph.PerturbWeights()
			graph.Forward(inA, inB)

			weights := make([]float64, len(graph.Synapses))
			for i, s := range graph.Synapses {
				weights[i] = s.Weight
			}

			packet := NeuralTelemetry{
				Timestamp:      time.Now().Format(time.RFC3339),
				Generation:     gen,
				FitnessScore:   graph.Fitness,
				NodeStates:     graph.Nodes,
				SynapseWeights: weights,
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Synaptic Evolution Router alive on ws://localhost:8080/neural-stream")
	_ = http.ListenAndServe(":8080", nil)
}
