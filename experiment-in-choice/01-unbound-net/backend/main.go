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

// DataPacket specifies our uniform JSON telemetry structure
type DataPacket struct {
	Timestamp   string  `json:"timestamp"`
	Heartbeat   int64   `json:"heartbeat"`
	Anomaly     bool    `json:"anomaly"`
	ResonanceHz float64 `json:"resonanceHz"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Allow local dashboard socket attachments
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[! ERR] Upgrade failure: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("[GO ENGINE] Low-latency connection mounted with system viewport client.")

	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	var heartbeat int64 = 0

	for {
		select {
		case <-ticker.C:
			heartbeat++
			isAnomaly := rand.Float64() < 0.08
			var resonance float64

			if isAnomaly {
				resonance = 32.0 + (rand.Float64() * 33.0)
			} else {
				resonance = 32.0 + (rand.Float64()*4.0 - 2.0)
			}

			packet := DataPacket{
				Timestamp:   time.Now().Format(time.RFC3339),
				Heartbeat:   heartbeat,
				Anomaly:     isAnomaly,
				ResonanceHz: math.Round(resonance*100) / 100,
			}

			message, err := json.Marshal(packet)
			if err != nil {
				fmt.Printf("[! ERR] Marshal error: %v\n", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("[GO ENGINE] Client channel detached. Flushing memory ring.")
				return
			}
		}
	}
}
func main() {
	http.HandleFunc("/", handleStream)
	fmt.Println("[INIT] High-performance Go networking matrix broadcasting on ws://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("[FATAL] Server launch crashed: %v\n", err)
	}
}
