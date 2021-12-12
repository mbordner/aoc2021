package main

import (
	"aoc2021/common/file"
	"fmt"
	"sort"
	"strings"
)

type node struct {
	id    string
	nodes map[string]*node
	big   bool
}

func newNode(id string) *node {
	n := new(node)
	n.id = id
	n.nodes = make(map[string]*node)
	if id[0] >= 'A' && id[0] <= 'Z' {
		n.big = true
	}
	return n
}

func (n *node) String() string {
	return n.id
}

func (n *node) add(o *node) {
	n.nodes[o.id] = o
}

type graph struct {
	start *node
	end   *node
	nodes map[string]*node
}

func newGraph() *graph {
	g := new(graph)
	g.nodes = make(map[string]*node)
	return g
}

func (g *graph) getNode(id string) *node {
	if n, ok := g.nodes[id]; ok {
		return n
	}
	n := newNode(id)
	g.add(n)
	return n
}

func (g *graph) add(n *node) {
	g.nodes[n.id] = n
}

func cloneVisited(visited map[string]int) map[string]int {
	m := make(map[string]int)
	for k, v := range visited {
		m[k] = v
	}
	return m
}

func markVisited(visited map[string]int, id string) {
	if v, ok := visited[id]; ok {
		visited[id] = v + 1
	} else {
		visited[id] = 1
	}
}

func getVisited(visited map[string]int, id string) int {
	if v, ok := visited[id]; ok {
		return v
	}
	return 0
}

func anyVisited(visited map[string]int, num int) bool {
	for _, v := range visited {
		if v == num {
			return true
		}
	}
	return false
}

func getPaths(n *node, visited map[string]int) [][]*node {
	paths := make([][]*node, 0, len(n.nodes)+1)
	if n.big == false {
		markVisited(visited, n.id)
	}
	if n.id == "end" {
		paths = append(paths, []*node{n})
	} else {
		for _, o := range n.nodes {
			newVisited := cloneVisited(visited)
			allowed := true
			if o.big == false {
				if o.id == "start" {
					allowed = false
				} else {
					numAllowed := 2
					if anyVisited(visited, 2) {
						numAllowed = 1
					}
					if getVisited(visited, o.id) > numAllowed-1 {
						allowed = false
					}
				}
			}
			if allowed {
				tpaths := getPaths(o, newVisited)
				if len(tpaths) > 0 {
					for _, tp := range tpaths {
						if tp[len(tp)-1].id == "end" {
							paths = append(paths, append([]*node{n}, tp...))
						}
					}
				}
			}

		}
	}
	return paths
}

func main() {
	g := getGraph("../data.txt")

	paths := getPaths(g.start, make(map[string]int))

	sort.Slice(paths, func(i, j int) bool {
		k := 0
		for ; k < len(paths[i]) && k < len(paths[j]) && paths[i][k].id == paths[j][k].id; k++ {
		}
		if paths[i][k] != paths[j][k] {
			return paths[i][k].id[0] < paths[j][k].id[0]
		}
		return len(paths[i]) < len(paths[k])
	})

	fmt.Println(len(paths))
}

func getGraph(filename string) *graph {
	lines, _ := file.GetLines(filename)

	g := newGraph()

	for _, line := range lines {
		tokens := strings.Split(line, "-")
		n := g.getNode(tokens[0])
		o := g.getNode(tokens[1])
		n.add(o)
		o.add(n)
	}

	g.start = g.getNode("start")
	g.end = g.getNode("end")

	return g
}
