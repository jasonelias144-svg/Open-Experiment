package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/cmplx"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// QuantumTelemetry packages circuit calculations for display viewports
type QuantumTelemetry struct {
	Timestamp     string    `json:"timestamp"`
	CycleStep     int64     `json:"cycleStep"`
	StateVector   []string  `json:"stateVector"`   // String representations of complex amplitudes
	Probabilities []float64 `json:"probabilities"` // Born rule probabilities: |ψ|^2
	Entanglement  float64   `json:"entanglement"`  // Simplified Von Neumann entropy indicator
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// QubitRegister holds state vectors for a 2-qubit space (dimension 4)
type QubitRegister struct {
	Amplitudes []complex128
}

func NewRegister() *QubitRegister {
	// Initialize to ground state |00>
	reg := &QubitRegister{
		Amplitudes: make([]complex128, 4),
	}
	reg.Amplitudes[0] = cmplx.Rect(1, 0) // Amplitude 1 for |00>
	return reg
}

// ApplyHadamard applies an H gate to the first qubit to trigger superposition
func (r *QubitRegister) ApplyHadamard() {
	invSqrt2 := 1.0 / math.Sqrt(2.0)
	hMat := complex(invSqrt2, 0)

	newAmps := make([]complex128, 4)
	// Matrix transformation mapping for Qubit 0
	newAmps[0] = hMat*r.Amplitudes[0] + hMat*r.Amplitudes[2]
	newAmps[1] = hMat*r.Amplitudes[1] + hMat*r.Amplitudes[3]
	newAmps[2] = hMat*r.Amplitudes[0] - hMat*r.Amplitudes[2]
	newAmps[3] = hMat*r.Amplitudes[1] - hMat*r.Amplitudes[3]
	r.Amplitudes = newAmps
}

// ApplyCNOT targets qubit 1 using qubit 0 as control to link entanglement vectors
func (r *QubitRegister) ApplyCNOT() {
	// Swaps amplitudes of |10> and |11> if control qubit is 1
	r.Amplitudes[2], r.Amplitudes[3] = r.Amplitudes[3], r.Amplitudes[2]
}
func (r *QubitRegister) GetMetrics() ([]string, []float64, float64) {
	ampsStr := make([]string, 4)
	probs := make([]float64, 4)

	for i, amp := range r.Amplitudes {
		ampsStr[i] = fmt.Sprintf("%.3f+%.3fi", real(amp), imag(amp))
		probs[i] = math.Pow(cmplx.Abs(amp), 2)
	}

	// Simplified structural measurement tracker for Bell State entanglement approximation
	// Measures how far state probabilities diverge from isolated baseline systems
	entanglementValue := 0.0
	if math.Abs(probs[0]-0.5) < 0.05 && math.Abs(probs[3]-0.5) < 0.05 {
		entanglementValue = 1.0 // State perfectly matches maximized Bell State |Φ+>
	}

	return ampsStr, probs, entanglementValue
}
func main() {
	http.HandleFunc("/quantum-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[QUANTUM] Observer socket attached to virtual register pipeline.")

		ticker := time.NewTicker(400 * time.Millisecond) // Coherent clock cycle gate
		defer ticker.Stop()

		var step int64 = 0

		for range ticker.C {
			step++
			reg := NewRegister()

			// Modulate quantum gate configurations over time blocks to create dynamic cycles
			phase := step % 4
			if phase >= 1 {
				reg.ApplyHadamard() // Enter superposition state
			}
			if phase >= 2 {
				reg.ApplyCNOT() // Trigger maximized qubit entanglement state
			}

			amps, probs, entang := reg.GetMetrics()

			packet := QuantumTelemetry{
				Timestamp:     time.Now().Format(time.RFC3339),
				CycleStep:     step,
				StateVector:   amps,
				Probabilities: probs,
				Entanglement:  entang,
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Quantum Simulator Core spinning on ws://localhost:8080/quantum-stream")
	_ = http.ListenAndServe(":8080", nil)
}
