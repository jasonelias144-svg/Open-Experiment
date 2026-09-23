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

// CodexTelemetry packages multidimensional plane values for display viewports
type CodexTelemetry struct {
	Timestamp      string             `json:"timestamp"`
	Generation     int64              `json:"generation"`
	ActivePlane    string             `json:"activePlane"`
	PlaneMetrics   map[string]float64 `json:"planeMetrics"`   // Real-time saturation: 0.0 to 1.0
	SystemDrift    float64            `json:"systemDrift"`    // Divergence from neutral root
	CurrentSymbols string             `json:"currentSymbols"` // Active glyph sequence
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// PlaneMatrix tracks state levels across the Seven Fields of Life
type PlaneMatrix struct {
	Physical      float64
	Psychological float64
	Information   float64
	Frequency     float64
	Astral        float64
	Spiritual     float64
	Universal     float64
}

func NewPlaneMatrix() *PlaneMatrix {
	return &PlaneMatrix{
		Physical:      0.5,
		Psychological: 0.5,
		Information:   0.5,
		Frequency:     0.5,
		Astral:        0.5,
		Spiritual:     0.5,
		Universal:     0.5,
	}
}
func (m *PlaneMatrix) Modulate(gen int64) (string, map[string]float64) {
	// Periodic wave functions representing Rhythm and Vibration cycles
	m.Physical = 0.5 + (0.3 * math.Sin(float64(gen)*0.1))
	m.Psychological = 0.5 + (0.25 * math.Cos(float64(gen)*0.08))
	m.Information = 0.5 + (0.2 * math.Sin(float64(gen)*0.15))

	// Simulate a sharp Anomaly/Spike (Cosmic Radiation or Trigger Event)
	if gen%12 == 0 {
		m.Frequency = 0.9
		m.Astral = rand.Float64()
	} else {
		m.Frequency = math.Max(0.1, m.Frequency*0.85)
		m.Astral = 0.5 + (0.1 * math.Sin(float64(gen)*0.05))
	}

	m.Spiritual = 0.5 + (0.35 * math.Sin(float64(gen)*0.03))
	m.Universal = 0.5 + (0.4 * math.Cos(float64(gen)*0.02))

	metrics := map[string]float64{
		"Physical":      m.Physical,
		"Psychological": m.Psychological,
		"Information":   m.Information,
		"Frequency":     m.Frequency,
		"Astral":        m.Astral,
		"Spiritual":     m.Spiritual,
		"Universal":     m.Universal,
	}

	// Determine the dominant active plane by finding the maximum value
	dominant := "Physical"
	maxVal := m.Physical
	for name, val := range metrics {
		if val > maxVal {
			maxVal = val
			dominant = name
		}
	}

	return dominant, metrics
}
func getGlyphSequence(plane string, anomaly bool) string {
	if anomaly {
		return "👁 → 🤲 → ❤️ → ⚡ → 🔄" // Full Return Loop
	}
	switch plane {
	case "Physical":
		return "👁 → 🤲 → ⚡"
	case "Psychological":
		return "👁 → ❤️ → 🪞"
	case "Information":
		return "🧭 → 🔍 → 💎"
	case "Frequency":
		return " VIBRATION // ⚡ × □"
	default:
		return "👁 → 🤲 → ❤️"
	}
}
func main() {
	matrix := NewPlaneMatrix()

	http.HandleFunc("/codex-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[CODEX] Observer client handshaked successfully with the Seven Planes Matrix.")

		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		var gen int64 = 0

		for range ticker.C {
			gen++
			dominantPlane, metrics := matrix.Modulate(gen)

			// Calculate drift using the divergence of Frequency from the baseline
			drift := math.Abs(metrics["Frequency"] - 0.5)

			packet := CodexTelemetry{
				Timestamp:      time.Now().Format(time.RFC3339),
				Generation:     gen,
				ActivePlane:    dominantPlane,
				PlaneMetrics:   metrics,
				SystemDrift:    drift,
				CurrentSymbols: getGlyphSequence(dominantPlane, drift > 0.35),
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Living Codex Engine online at ws://localhost:8080/codex-stream")
	_ = http.ListenAndServe(":8080", nil)
}
