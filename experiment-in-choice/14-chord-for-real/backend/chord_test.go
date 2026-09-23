package main

import (
	"math"
	"testing"
)

func TestBetween(t *testing.T) {
	r := NewRing(16, 1)
	cases := []struct {
		x, a, b int
		want    bool
	}{
		{5, 3, 8, true}, {3, 3, 8, false}, {8, 3, 8, false},
		{15, 12, 2, true}, {1, 12, 2, true}, {5, 12, 2, false},
		{7, 4, 4, true}, {4, 4, 4, false}, // a == b: everything except a
	}
	for _, c := range cases {
		if got := r.between(c.x, c.a, c.b); got != c.want {
			t.Errorf("between(%d, %d, %d) = %v, want %v", c.x, c.a, c.b, got, c.want)
		}
	}
}

func TestStableRingRoutesCorrectlyInLogHops(t *testing.T) {
	for _, n := range []int{1, 2, 3, 64, 512} {
		r, rounds := BuildStable(n, 1<<16, int64(n))
		if !r.Converged() {
			t.Fatalf("N=%d did not converge (%d rounds)", n, rounds)
		}
		hops, acc := r.LookupStats(5000)
		if acc != 1 {
			t.Errorf("N=%d: stable ring answered %.4f of lookups correctly, want all", n, acc)
		}
		if n > 1 {
			if bound := math.Log2(float64(n)) + 1; hops > bound {
				t.Errorf("N=%d: mean hops %.2f exceeds log2 N + 1 = %.2f", n, hops, bound)
			}
		}
	}
}

func TestRingRecoversAfterChurn(t *testing.T) {
	r, _ := BuildStable(128, 1<<16, 3)
	for i := 0; i < 40; i++ { // replace nearly a third of the ring with no maintenance in between
		r.Leave(r.RandomNode())
		r.Join(r.FreeID())
	}
	rounds := 0
	for !r.Converged() && rounds < 1000 {
		r.Round()
		rounds++
	}
	if !r.Converged() {
		t.Fatalf("ring did not repair itself after churn")
	}
	if _, acc := r.LookupStats(2000); acc != 1 {
		t.Errorf("after repair, accuracy %.4f, want 1", acc)
	}
	t.Logf("repaired in %d rounds", rounds)
}
