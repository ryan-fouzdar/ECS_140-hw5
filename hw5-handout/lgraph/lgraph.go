package lgraph

import (
	"sync"
)

type node uint

type edge struct {
	destination node
	label       rune
}

// LGraph is a function representing a directed labeled graph. If the node exists
// in the graph, the function returns true along with the set of outgoing edges
// from that node, otherwise false and nil.
type LGraph func(node) ([]edge, bool)

func check(g LGraph, s, t node, seq []rune) bool {
	edges, exists := g(s)
	if !exists {
		return false
	}
	if len(seq) == 0 {
		return s==t
	}
	for _, e := range edges {
		if e.label == seq[0] {
			if reached := check(g, e.destination, t, seq[1:]); reached {
				return true
			}
		}
	}
	return false
}

func find(g1, g2 LGraph, c, t node, k uint, s node, prefix []rune, result chan<- []rune, wg *sync.WaitGroup, mu *sync.Mutex, found *bool) {
	defer wg.Done()

	mu.Lock()
	if *found {
		mu.Unlock()
		return
	}
	mu.Unlock()

	edges, exists := g1(c)
	if !exists {
		return
	}

	if k == 0 {
		if c == t && !check(g2, s, t, prefix) {
			mu.Lock()
			if !*found {
				*found = true
				result <- prefix
			}
			mu.Unlock()
		}
		return
	}

	for _, e := range edges {
		newPrefix := append([]rune{}, prefix...) 
		newPrefix = append(newPrefix, e.label)

		wg.Add(1)
		go find(g1, g2, e.destination, t, k-1, s, newPrefix, result, wg, mu, found)
	}
}

// FindSequence searches for a sequence in g1 that doesn't exist in g2 concurrently.
func FindSequence(g1, g2 LGraph, s, t node, k uint) ([]rune, bool) {
	var wg sync.WaitGroup
	result := make(chan []rune, 1) 
	var mu sync.Mutex
	found := false

	wg.Add(1)
	go find(g1, g2, s, t, k, s, []rune{}, result, &wg, &mu, &found)

	go func() {
		wg.Wait()
		close(result)
	}()

	select {
	case seq, ok := <-result:
		if ok {
			return seq, true
		}
		return nil, false
	}
}
