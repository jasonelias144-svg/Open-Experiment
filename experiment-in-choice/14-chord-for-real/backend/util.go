package main

import (
	"crypto/sha1"
	"encoding/binary"
	"math"
)

// hashKey places a key on the ring with SHA-1, as Chord does.
func hashKey(k string, m int) int {
	sum := sha1.Sum([]byte(k))
	return int(binary.BigEndian.Uint32(sum[:4]) % uint32(m))
}

// loadSpread is the standard deviation of keys per node (cycle 13 called it "entropy").
func loadSpread(load map[int]int, nodes, keys int) float64 {
	if nodes == 0 {
		return 0
	}
	mean := float64(keys) / float64(nodes)
	sq := 0.0
	for _, c := range load {
		d := float64(c) - mean
		sq += d * d
	}
	sq += float64(nodes-len(load)) * mean * mean // nodes holding no keys
	return math.Sqrt(sq / float64(nodes))
}
