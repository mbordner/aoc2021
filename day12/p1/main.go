package main

import (
	"aoc2021/common/file"
	"fmt"
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

func cloneVisited(visited map[string]bool) map[string]bool {
	m := make(map[string]bool)
	for k, v := range visited {
		m[k] = v
	}
	return m
}

func getPaths(n *node, visited map[string]bool) [][]*node {
	paths := make([][]*node, 0, len(n.nodes)+1)
	if n.big == false {
		visited[n.id] = true
	}
	if n.id == "end" {
		paths = append(paths, []*node{n})
	} else {
		for _, o := range n.nodes {
			newVisited := cloneVisited(visited)
			allowed := true
			if o.big == false {
				if _, exists := newVisited[o.id]; !exists {
					newVisited[o.id] = true
				} else {
					allowed = false
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

	paths := getPaths(g.start, make(map[string]bool))

	fmt.Println(len(paths), paths)
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
