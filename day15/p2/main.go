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
	_, cost := shortestPaths.GetShortestPath(c.goal)

	fmt.Println(cost)

}

func getValue(oval int, inc int) int {
	val := ((oval-1)+inc)%9 + 1
	return val
}

func getCave(filename string) *cave {
	c := new(cave)
	c.g = graph.NewGraph()
	c.bb = &geom.BoundingBox{}

	lines, _ := file.GetLines(filename)

	ovalues := make([][]int, len(lines), len(lines))
	for y, line := range lines {
		ovalues[y] = make([]int, len(line), len(line))
		for x, char := range line {
			v, _ := strconv.ParseInt(string(char), 10, 32)
			ovalues[y][x] = int(v)

		}
	}

	oysize := len(ovalues)
	oxsize := len(ovalues[0])

	ysize := len(lines) * 5
	xsize := len(lines[0]) * 5

	values := make([][]int, ysize, ysize)
	for y := range values {
		values[y] = make([]int, xsize, xsize)
	}

	tiles := 5

	for i := 0; i < tiles; i++ {
		for j := 0; j < tiles; j++ {

			for y := 0; y < oysize; y++ {
				for x := 0; x < oxsize; x++ {
					values[y+(i*oysize)][x+(j*oxsize)] = getValue(ovalues[y][x], i+j)
				}
			}

		}
	}

	for y := range values {
		for x := range values {
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
			} else if x == xsize-1 && y == ysize-1 {
				c.goal = n
			}

			var e *graph.Edge

			if y > 0 {
				above := c.g.GetNode(geom.Pos{Y: y - 1, X: x})
				e = n.AddEdge(above, float64(above.GetProperty("value").(int)))
				e.AddProperty("dir", geom.North)
			}

			if y < ysize-1 {
				below := c.g.GetNode(geom.Pos{Y: y + 1, X: x})
				e = n.AddEdge(below, float64(below.GetProperty("value").(int)))
				e.AddProperty("dir", geom.South)
			}

			if x > 0 {
				left := c.g.GetNode(geom.Pos{Y: y, X: x - 1})
				e = n.AddEdge(left, float64(left.GetProperty("value").(int)))
				e.AddProperty("dir", geom.West)
			}

			if x < xsize-1 {
				right := c.g.GetNode(geom.Pos{Y: y, X: x + 1})
				e = n.AddEdge(right, float64(right.GetProperty("value").(int)))
				e.AddProperty("dir", geom.East)
			}

		}
	}

	return c
}
