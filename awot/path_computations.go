package awot

import (
	"gonum.org/v1/gonum/graph"
	"log"
)

type Path []graph.Node

func (p Path) toArray() []graph.Node {
	return []graph.Node(p)
}

// Return all the intersections of any t paths in given paths
func comb(paths []Path, t int) []Path {

	if t == len(paths) {
		return paths
	}
	if t == 0 {
		return []Path{}
	}

	with := comb(paths[1:], t-1)
	for i, _ := range with {
		with[i] = intersection([][]graph.Node{paths[0], with[i]})
	}

	without := comb(paths[1:], t)

	return append(with, without...)
}

// Compute the probability of the given shortest paths, using the inclusion exclusion formula
func (ring KeyRing) probabilityOfMinPaths(minpaths [][]graph.Node) float32 {
  // convert minpaths to []Path
  minPaths := make([]Path, len(minpaths))
  for i, v := range minpaths {
    minPaths[i] = Path(v)
  }

	p := float32(0.0)

	n := len(minPaths)
	s := float32(1.0)
	for i := 0; i < n; i++ {
		// get the possible paths of intersection of i paths in the n given
		// n choose i such paths
		npaths := comb(minPaths, i)
		for _, path := range npaths {
			p += s * ring.probabilityOfPath(path)
		}

		s = -s
	}

  return p
}

// Compute the probability of the given path
func (ring KeyRing) probabilityOfPath(path []graph.Node) float32 {
	p := float32(1.0)

	for i, node := range path {
		if i == 0 || i == len(path) {
			// do not include source nor terminal
			continue
		}

		v, present := ring.getVertex(node)
		if !present {
			p = 0
			log.Fatal("NODE DOES NOT EXISTS in ids !")
		} else {
			p = p * *v.probability
		}
	}

	return p
}

// Return the intersection of given paths
// e.g. A={1,2,3} B={2,4,5}, A inter B = {1,2,3,4,5}
func intersection(paths [][]graph.Node) []graph.Node {
	var nodes []graph.Node
	for _, path := range paths {
		nodes = append(nodes, path...)
	}

	encountered := make(map[int64]bool, 0)
	r := make([]graph.Node, 0)

	for _, n := range nodes {
		if !encountered[n.ID()] {
			// not yet added
			r = append(r, n)
			encountered[n.ID()] = true
		}
	}
	return r
}
