package lgraph

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
func find(g1, g2 LGraph, c, t node, k uint, s node, prefix []rune) ([]rune, bool) {
	edges, exists := g1(c)

	if !exists {
		return nil, false
	}

	if k ==0 {
		if c == t && !check(g2, s, t, prefix) {
			return prefix, true
		}
		return nil, false
	}

	for _, e := range edges {
		if seq, found := find(g1, g2, e.destination, t, k-1, s, append(prefix, e.label)); found {
			return seq, found
		}
	}
	return nil, false
}
// FindSequence returns (S, true) if there is a sequence S of length k from node
// s to node t in graph g1 and S is not a sequence from s to t in graph g2; else
// it returns (nil, false).
func FindSequence(g1, g2 LGraph, s, t node, k uint) ([]rune, bool) {
	return find(g1, g2, s, t, k, s, []rune{})
}