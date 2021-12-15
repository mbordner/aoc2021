package main

import (
	"aoc2021/common/file"
	"aoc2021/common/geom"
	"aoc2021/common/graph"
	"aoc2021/common/graph/djikstra"
	"fmt"
	"strconv"
)

type cave struct {
	g     *graph.Graph
	bb    *geom.BoundingBox
	start *graph.Node
	goal  *graph.Node
}

func main() {
	c := getCave("../data.txt")

	shortestPaths := djikstra.GenerateShortestPaths(c.g, c.start)
	shortestPath, cost := shortestPaths.GetShortestPath(c.goal)

	fmt.Println(shortestPath, cost)

}

func getCave(filename string) *cave {
	c := new(cave)
	c.g = graph.NewGraph()
	c.bb = &geom.BoundingBox{}

	lines, _ := file.GetLines(filename)

	values := make([][]int, len(lines), len(lines))
	for y, line := range lines {
		values[y] = make([]int, len(line), len(line))
		for x, char := range line {
			v, _ := strconv.ParseInt(string(char), 10, 32)
			values[y][x] = int(v)

			point := geom.Pos{Y: y, X: x}
			c.bb.Extend(point)

			n := c.g.CreateNode(point)
			n.AddProperty("value", values[y][x])
		}
	}

	for y := range values {
		for x := range values {

			n := c.g.GetNode(geom.Pos{Y: y, X: x})

			if x == 0 && y == 0 {
				c.start = n
			} else if x == len(lines)-1 && y == len(lines[y])-1 {
				c.goal = n
			}

			var e *graph.Edge

			if y > 0 {
				above := c.g.GetNode(geom.Pos{Y: y - 1, X: x})
				e = n.AddEdge(above, float64(above.GetProperty("value").(int)))
				e.AddProperty("dir", geom.North)
			}

			if y < len(values)-1 {
				below := c.g.GetNode(geom.Pos{Y: y + 1, X: x})
				e = n.AddEdge(below, float64(below.GetProperty("value").(int)))
				e.AddProperty("dir", geom.South)
			}

			if x > 0 {
				left := c.g.GetNode(geom.Pos{Y: y, X: x - 1})
				e = n.AddEdge(left, float64(left.GetProperty("value").(int)))
				e.AddProperty("dir", geom.West)
			}

			if x < len(values[y])-1 {
				right := c.g.GetNode(geom.Pos{Y: y, X: x + 1})
				e = n.AddEdge(right, float64(right.GetProperty("value").(int)))
				e.AddProperty("dir", geom.East)
			}

		}
	}

	return c
}
