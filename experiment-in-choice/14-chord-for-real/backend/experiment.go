package main

import (
	"fmt"
	"math"
)

// LookupStats routes `samples` random lookups from random nodes and reports the
// mean hop count and the fraction that reached the true owner.
func (r *Ring) LookupStats(samples int) (meanHops, accuracy float64) {
	s := r.sorted()
	hops, correct := 0, 0
	for i := 0; i < samples; i++ {
		start := s[r.rng.Intn(len(s))]
		key := r.rng.Intn(r.M)
		got, h := r.FindSuccessor(start, key)
		hops += h
		if got == trueSuccessor(s, key) {
			correct++
		}
	}
	return float64(hops) / float64(samples), float64(correct) / float64(samples)
}

// BuildStable joins n nodes one at a time (one maintenance round after each
// join) and then runs rounds until every pointer is correct. It returns the
// ring and the number of extra rounds convergence took.
func BuildStable(n, m int, seed int64) (*Ring, int) {
	r := NewRing(m, seed)
	for i := 0; i < n; i++ {
		r.Join(r.FreeID())
		r.Round()
	}
	rounds := 0
	for !r.Converged() && rounds < 1000 {
		r.Round()
		rounds++
	}
	return r, rounds
}

// RunExperiment prints the two results in the README.
func RunExperiment() {
	const m = 1 << 20
	fmt.Println("1) Stable ring: mean lookup hops vs the Chord paper's ½·log2(N)")
	fmt.Println()
	fmt.Println("| N | rounds to converge | mean hops | ½·log2 N | ratio | lookups correct |")
	fmt.Println("|---:|---:|---:|---:|---:|---:|")
	for n := 8; n <= 2048; n *= 2 {
		r, rounds := BuildStable(n, m, int64(n))
		hops, acc := r.LookupStats(20000)
		theory := 0.5 * math.Log2(float64(n))
		fmt.Printf("| %d | %d | %.2f | %.2f | %.2f | %.2f%% |\n", n, rounds, hops, theory, hops/theory, 100*acc)
	}

	fmt.Println()
	fmt.Println("2) Under churn: N = 256, each round a fraction of nodes leave and as many join")
	fmt.Println()
	fmt.Println("| nodes replaced per round | stale fingers | mean hops | lookups correct |")
	fmt.Println("|---:|---:|---:|---:|")
	for _, rate := range []float64{0, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2} {
		r, _ := BuildStable(256, m, 7)
		var hopSum, accSum, staleSum float64
		const warmup, measure = 50, 200
		for round := 0; round < warmup+measure; round++ {
			// Fractional rates are applied as an expected count per round.
			k := int(rate * 256)
			if r.rng.Float64() < rate*256-float64(k) {
				k++
			}
			for j := 0; j < k; j++ {
				r.Leave(r.RandomNode())
				r.Join(r.FreeID())
			}
			r.Round()
			if round >= warmup {
				h, a := r.LookupStats(200)
				stale, total := r.StaleFingers()
				hopSum += h
				accSum += a
				staleSum += float64(stale) / float64(total)
			}
		}
		fmt.Printf("| %.1f%% | %.1f%% | %.2f | %.2f%% |\n", 100*rate, 100*staleSum/measure, hopSum/measure, 100*accSum/measure)
	}
}
