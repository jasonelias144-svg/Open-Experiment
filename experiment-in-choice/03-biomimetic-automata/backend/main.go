package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	GridWidth  = 64
	GridHeight = 64
)

// CellPacket packages cellular state logs for the visualization layer
type CellPacket struct {
	Timestamp       string `json:"timestamp"`
	Generation      int64  `json:"generation"`
	ActiveCellCount int    `json:"activeCount"`
	GridState       []int  `json:"gridState"` // Flattened 1D array binary snapshot
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Global grid matrices
type Biosphere struct {
	Current [GridWidth][GridHeight]int
	Next    [GridWidth][GridHeight]int
}

func NewBiosphere() *Biosphere {
	b := &Biosphere{}
	// Seed with a randomized structural matrix baseline
	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			if rand.Float64() < 0.25 {
				b.Current[x][y] = 1
			}
		}
	}
	return b
}
func (b *Biosphere) Step() int {
	activeCount := 0
	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			neighbors := b.countNeighbors(x, y)
			alive := b.Current[x][y] == 1

			// Standard Conway Automata parameters: Birth, Survival, and Overpopulation Decay
			if alive && (neighbors < 2 || neighbors > 3) {
				b.Next[x][y] = 0
			} else if !alive && neighbors == 3 {
				b.Next[x][y] = 1
			} else {
				b.Next[x][y] = b.Current[x][y]
			}

			if b.Next[x][y] == 1 {
				activeCount++
			}
		}
	}
	b.Current = b.Next
	return activeCount
}
func (b *Biosphere) countNeighbors(cx, cy int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			// Wrap bounds mapping to simulate an infinite toroidal universe
			x := (cx + i + GridWidth) % GridWidth
			y := (cy + j + GridHeight) % GridHeight
			if b.Current[x][y] == 1 {
				count++
			}
		}
	}
	return count
}
func (b *Biosphere) Flatten() []int {
	flat := make([]int, GridWidth*GridHeight)
	idx := 0
	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			flat[idx] = b.Current[x][y]
			idx++
		}
	}
	return flat
}
func main() {
	biosphere := NewBiosphere()

	http.HandleFunc("/bio-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[BIOSPHERE] Client attached to metabolic cell stream pipeline.")

		ticker := time.NewTicker(150 * time.Millisecond) // Core heartbeat tick
		defer ticker.Stop()

		var gen int64 = 0

		for range ticker.C {
			gen++
			activeCells := biosphere.Step()

			packet := CellPacket{
				Timestamp:       time.Now().Format(time.RFC3339),
				Generation:      gen,
				ActiveCellCount: activeCells,
				GridState:       biosphere.Flatten(),
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Cellular Automata Engine online at ws://localhost:8080/bio-stream")
	_ = http.ListenAndServe(":8080", nil)
}
