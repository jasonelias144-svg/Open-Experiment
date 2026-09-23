package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ChordTelemetry keeps the exact fields cycle 13's page reads, so that page
// works unchanged. routingSteps is now a measured hop count, not a random one.
type ChordTelemetry struct {
	Timestamp      string           `json:"timestamp"`
	LifecycleTick  int64            `json:"lifecycleTick"`
	NodePositions  []int            `json:"nodePositions"`
	KeyAssignments map[int][]string `json:"keyAssignments"`
	RoutingSteps   int              `json:"routingSteps"`
	NetworkEntropy float64          `json:"networkEntropy"`
	// Added in this cycle.
	LookupCorrect bool    `json:"lookupCorrect"`
	StaleFingers  float64 `json:"staleFingers"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func main() {
	experiment := flag.Bool("experiment", false, "print the hop-count and churn tables, then exit")
	flag.Parse()
	if *experiment {
		RunExperiment()
		return
	}

	// Ring size 1000 so positions fit the 0–999 layout of cycle 13's page.
	ring, _ := BuildStable(10, 1000, time.Now().UnixNano())
	var keys []string
	var mu sync.Mutex
	var tick int64

	http.HandleFunc("/chord-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			mu.Lock()
			tick++
			switch {
			case tick%9 == 0 && len(ring.Nodes) > 5:
				ring.Leave(ring.RandomNode())
			case tick%7 == 0 && len(ring.Nodes) < 16:
				ring.Join(ring.FreeID())
			}
			if tick%3 == 0 {
				keys = append(keys, fmt.Sprintf("key_%04d", tick))
				if len(keys) > 40 {
					keys = keys[1:]
				}
			}
			ring.Round()

			// One real lookup per frame, from a random node, for a random key.
			s := ring.sorted()
			target := ring.rng.Intn(ring.M)
			got, hops := ring.FindSuccessor(s[ring.rng.Intn(len(s))], target)

			positions := make([]int, len(s))
			assign := map[int][]string{}
			load := map[int]int{}
			for i, n := range s {
				positions[i] = n.ID
			}
			for _, k := range keys {
				owner := trueSuccessor(s, hashKey(k, ring.M)).ID
				assign[owner] = append(assign[owner], k)
				load[owner]++
			}
			stale, total := ring.StaleFingers()
			packet := ChordTelemetry{
				Timestamp:      time.Now().Format(time.RFC3339),
				LifecycleTick:  tick,
				NodePositions:  positions,
				KeyAssignments: assign,
				RoutingSteps:   hops,
				NetworkEntropy: loadSpread(load, len(s), len(keys)),
				LookupCorrect:  got == trueSuccessor(s, target),
				StaleFingers:   float64(stale) / float64(total),
			}
			mu.Unlock()
			msg, _ := json.Marshal(packet)
			if conn.WriteMessage(websocket.TextMessage, msg) != nil {
				return
			}
		}
	})
	fmt.Println("[CHORD] real finger-table routing on ws://localhost:8080/chord-stream")
	_ = http.ListenAndServe(":8080", nil)
}
