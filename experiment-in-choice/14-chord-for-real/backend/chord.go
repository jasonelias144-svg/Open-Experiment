package main

import (
	"math/rand"
	"sort"
)

// Node is one Chord peer. It only knows its successor, its predecessor and its
// finger table. Everything else has to be discovered by routing.
type Node struct {
	ID          int
	Successor   *Node
	Predecessor *Node
	Fingers     []*Node // Fingers[i] should be successor(ID + 2^i)
	next        int     // next finger to refresh (fix_fingers is round-robin)
	alive       bool
}

// Ring is a simulated Chord network on an identifier circle of size M.
// The simulator itself can see every node (for ground truth and statistics),
// but routing only ever follows the pointers each node holds.
type Ring struct {
	M     int
	Nodes map[int]*Node
	rng   *rand.Rand
	bits  int // number of fingers: smallest b with 2^b >= M
}

func NewRing(m int, seed int64) *Ring {
	b := 0
	for 1<<b < m {
		b++
	}
	return &Ring{M: m, Nodes: map[int]*Node{}, rng: rand.New(rand.NewSource(seed)), bits: b}
}

// between reports whether x lies on the open arc (a, b), going clockwise.
// When a == b the arc is the whole circle except a.
func (r *Ring) between(x, a, b int) bool {
	if a < b {
		return a < x && x < b
	}
	return x > a || x < b
}

// betweenRightIncl reports whether x lies on the half-open arc (a, b].
func (r *Ring) betweenRightIncl(x, a, b int) bool {
	return x == b || r.between(x, a, b)
}

func (r *Ring) closestPreceding(n *Node, id int) *Node {
	for i := len(n.Fingers) - 1; i >= 0; i-- {
		f := n.Fingers[i]
		if f != nil && f.alive && r.between(f.ID, n.ID, id) {
			return f
		}
	}
	return n
}

// FindSuccessor routes a lookup for id starting at node n, the way a Chord peer
// would: forward to the closest preceding finger until id falls between the
// current node and its successor. It returns the answer and the number of
// node-to-node hops the query took.
func (r *Ring) FindSuccessor(n *Node, id int) (*Node, int) {
	hops := 0
	for !r.betweenRightIncl(id, n.ID, n.Successor.ID) {
		next := r.closestPreceding(n, id)
		if next == n {
			next = n.Successor // no useful finger yet: fall back to walking the ring
		}
		n = next
		hops++
		if hops > r.M { // only possible if the successor ring itself is broken
			break
		}
	}
	return n.Successor, hops
}

// Join adds a node with the given id, bootstrapping through a random live node.
// Like real Chord it only sets the successor; stabilization does the rest.
func (r *Ring) Join(id int) *Node {
	n := &Node{ID: id, Fingers: make([]*Node, r.bits), alive: true}
	if len(r.Nodes) == 0 {
		n.Successor = n
	} else {
		boot := r.RandomNode()
		n.Successor, _ = r.FindSuccessor(boot, id)
	}
	r.Nodes[id] = n
	return n
}

// Leave removes a node gracefully: it hands its neighbours to each other, as in
// the Chord paper's voluntary departure. Other nodes' fingers may still point
// at it until fix_fingers replaces them. Routing skips dead fingers meanwhile.
func (r *Ring) Leave(n *Node) {
	if len(r.Nodes) <= 1 {
		return
	}
	for _, p := range r.Nodes { // whoever has n as successor is its predecessor
		if p.Successor == n && p != n {
			p.Successor = n.Successor
		}
	}
	if n.Successor.Predecessor == n {
		n.Successor.Predecessor = n.Predecessor
	}
	n.alive = false
	delete(r.Nodes, n.ID)
}

// Stabilize is the periodic check that keeps successor pointers correct.
func (r *Ring) Stabilize(n *Node) {
	if x := n.Successor.Predecessor; x != nil && x.alive && r.between(x.ID, n.ID, n.Successor.ID) {
		n.Successor = x
	}
	s := n.Successor
	if s.Predecessor == nil || !s.Predecessor.alive || r.between(n.ID, s.Predecessor.ID, s.ID) {
		s.Predecessor = n
	}
}

// FixNextFinger refreshes one finger per call, round-robin.
func (r *Ring) FixNextFinger(n *Node) {
	n.next = (n.next + 1) % r.bits
	n.Fingers[n.next], _ = r.FindSuccessor(n, (n.ID+1<<n.next)%r.M)
}

// Round runs one maintenance period: every node stabilizes and fixes one finger.
func (r *Ring) Round() {
	for _, n := range r.sorted() {
		r.Stabilize(n)
		r.FixNextFinger(n)
	}
}

// ---- ground truth, visible only to the simulator ----

func (r *Ring) sorted() []*Node {
	out := make([]*Node, 0, len(r.Nodes))
	for _, n := range r.Nodes {
		out = append(out, n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// TrueSuccessor is the node that really owns id: the first live node at or
// clockwise after it.
func (r *Ring) TrueSuccessor(id int) *Node { return trueSuccessor(r.sorted(), id) }

func trueSuccessor(s []*Node, id int) *Node {
	i := sort.Search(len(s), func(i int) bool { return s[i].ID >= id })
	return s[i%len(s)]
}

// StaleFingers counts fingers that differ from what a fully repaired ring would hold.
func (r *Ring) StaleFingers() (stale, total int) {
	s := r.sorted()
	for _, n := range s {
		for i, f := range n.Fingers {
			total++
			if f != trueSuccessor(s, (n.ID+1<<i)%r.M) {
				stale++
			}
		}
	}
	return
}

// Converged reports whether every successor pointer and finger is correct.
func (r *Ring) Converged() bool {
	s := r.sorted()
	for _, n := range s {
		if n.Successor != trueSuccessor(s, (n.ID+1)%r.M) {
			return false
		}
	}
	stale, _ := r.StaleFingers()
	return stale == 0
}

func (r *Ring) RandomNode() *Node {
	s := r.sorted()
	return s[r.rng.Intn(len(s))]
}

// FreeID picks an identifier no live node is using.
func (r *Ring) FreeID() int {
	for {
		id := r.rng.Intn(r.M)
		if _, taken := r.Nodes[id]; !taken {
			return id
		}
	}
}
