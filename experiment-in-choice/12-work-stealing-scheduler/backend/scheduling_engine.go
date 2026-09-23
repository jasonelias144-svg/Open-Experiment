package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// SchedulerTelemetry packages load balancing statistics for the visualization panel
type SchedulerTelemetry struct {
	Timestamp       string  `json:"timestamp"`
	SchedulerCycle  int64   `json:"schedulerCycle"`
	QueueSizes      []int   `json:"queueSizes"`      // Remaining computational tasks per worker
	StealEvents     int64   `json:"stealEvents"`     // Cumulative work-stealing occurrences
	TotalThroughput int64   `json:"totalThroughput"` // Successfully resolved processing blocks
	SystemBalance   float64 `json:"systemBalance"`   // Gini coefficient of thread allocation efficiency
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WorkerActor struct {
	ID        int
	TaskQueue []int // Slice acting as a basic work queue deck
	Mu        sync.Mutex
}
type CentralScheduler struct {
	Workers    []*WorkerActor
	Steals     int64
	Throughput int64
	CycleCount int64
	Mu         sync.Mutex
}

func NewCentralScheduler(workerCount int) *CentralScheduler {
	s := &CentralScheduler{
		Workers: make([]*WorkerActor, workerCount),
		Steals:  0,
	}
	for i := 0; i < workerCount; i++ {
		s.Workers[i] = &WorkerActor{ID: i, TaskQueue: []int{}}
		// Initially saturate only Worker 0 to create severe load imbalances
		if i == 0 {
			for t := 0; t < 30; t++ {
				s.Workers[i].TaskQueue = append(s.Workers[i].TaskQueue, rand.Intn(10)+1)
			}
		}
	}
	return s
}
func (s *CentralScheduler) ExecuteWorkCycle() {
	s.Mu.Lock()
	s.CycleCount++
	s.Mu.Unlock()

	// Inject periodic fresh incoming task storms directly into random queues
	if s.CycleCount%8 == 0 {
		target := rand.Intn(len(s.Workers))
		s.Workers[target].Mu.Lock()
		for t := 0; t < 12; t++ {
			s.Workers[target].TaskQueue = append(s.Workers[target].TaskQueue, rand.Intn(5)+1)
		}
		s.Workers[target].Mu.Unlock()
	}

	// Simulating parallel independent worker actor execution sweeps
	for _, worker := range s.Workers {
		worker.Mu.Lock()

		// If local queue contains computing blocks, process the topmost element
		if len(worker.TaskQueue) > 0 {
			worker.TaskQueue = worker.TaskQueue[1:]
			worker.Mu.Unlock()
			s.Mu.Lock()
			s.Throughput++
			s.Mu.Unlock()
		} else {
			// Local queue empty: Enter Work-Stealing Protocol mode
			worker.Mu.Unlock()
			s.attemptWorkSteal(worker.ID)
		}
	}
}
func (s *CentralScheduler) attemptWorkSteal(stealerID int) {
	// Pick a random victim thread actor to rob code blocks from
	victimID := rand.Intn(len(s.Workers))
	if victimID == stealerID {
		return
	}

	victim := s.Workers[victimID]
	stealer := s.Workers[stealerID]

	victim.Mu.Lock()
	// Victim must contain transferable work items to trigger execution steal
	if len(victim.TaskQueue) > 4 {
		stealSegmentSize := len(victim.TaskQueue) / 2
		stolenTasks := victim.TaskQueue[len(victim.TaskQueue)-stealSegmentSize:]
		victim.TaskQueue = victim.TaskQueue[:len(victim.TaskQueue)-stealSegmentSize]
		victim.Mu.Unlock()

		stealer.Mu.Lock()
		stealer.TaskQueue = append(stealer.TaskQueue, stolenTasks...)
		stealer.Mu.Unlock()

		s.Mu.Lock()
		s.Steals++
		s.Mu.Unlock()
		fmt.Printf("[STEAL] Actor_0%d successfully poached %d tasks from Actor_0%d\n", stealerID, stealSegmentSize, victimID)
	} else {
		victim.Mu.Unlock()
	}
}
func (s *CentralScheduler) CalculateGiniEfficiency() float64 {
	// Computes statistical structural balance deviation across worker pools
	var sizes []float64
	var sum float64
	for _, w := range s.Workers {
		w.Mu.Lock()
		sz := float64(len(w.TaskQueue))
		sizes = append(sizes, sz)
		sum += sz
		w.Mu.Unlock()
	}
	if sum == 0 {
		return 1.0
	}
	var absoluteDiff float64
	for _, xi := range sizes {
		for _, xj := range sizes {
			absoluteDiff += math.Abs(xi - xj)
		}
	}
	n := float64(len(s.Workers))
	return 1.0 - (absoluteDiff / (2.0 * n * sum))
}
func main() {
	scheduler := NewCentralScheduler(4) // 4-Thread Symmetric Processing Architecture

	http.HandleFunc("/scheduler-stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Println("[SCHEDULER] Low-latency metric observer linked onto work-stealing core core.")

		ticker := time.NewTicker(100 * time.Millisecond) // Engine execution speed gate
		defer ticker.Stop()

		for range ticker.C {
			scheduler.ExecuteWorkCycle()
			balanceMetric := scheduler.CalculateGiniEfficiency()

			scheduler.Mu.Lock()
			qSizes := make([]int, len(scheduler.Workers))
			for i, w := range scheduler.Workers {
				w.Mu.Lock()
				qSizes[i] = len(w.TaskQueue)
				w.Mu.Unlock()
			}
			cycle := scheduler.CycleCount
			steals := scheduler.Steals
			throughput := scheduler.Throughput
			scheduler.Mu.Unlock()

			packet := SchedulerTelemetry{
				Timestamp:       time.Now().Format(time.RFC3339),
				SchedulerCycle:  cycle,
				QueueSizes:      qSizes,
				StealEvents:     steals,
				TotalThroughput: throughput,
				SystemBalance:   balanceMetric,
			}

			payload, _ := json.Marshal(packet)
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		}
	})

	fmt.Println("[INIT] Distributed Work-Stealing Core listening on ws://localhost:8080/scheduler-stream")
	_ = http.ListenAndServe(":8080", nil)
}
