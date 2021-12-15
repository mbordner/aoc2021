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
	fmt.Println("graph constructed")

	shortestPaths := djikstra.GenerateShortestPaths(c.g, c.start)
	fmt.Println("shortest paths calculated from start")

	_, cost := shortestPaths.GetShortestPath(c.goal)

	fmt.Println(cost)

}

func getIncrementingTileValue(oval int, inc int) int {
	val := ((oval-1)+inc)%9 + 1
	return val
}

func getCave(filename string) *cave {
	c := new(cave)
	c.g = graph.NewGraph()
	c.bb = &geom.BoundingBox{}

	lines, _ := file.GetLines(filename)

	origValues := make([][]int, len(lines), len(lines))
	for y, line := range lines {
		origValues[y] = make([]int, len(line), len(line))
		for x, char := range line {
			v, _ := strconv.ParseInt(string(char), 10, 32)
			origValues[y][x] = int(v)

		}
	}

	numTiles := 5

	origYDimension := len(origValues)
	origXDimension := len(origValues[0])

	yDimension := len(lines) * numTiles
	xDimension := len(lines[0]) * numTiles

	values := make([][]int, yDimension, yDimension)
	for y := range values {
		values[y] = make([]int, xDimension, xDimension)
	}

	for i := 0; i < numTiles; i++ {
		for j := 0; j < numTiles; j++ {

			for y := 0; y < origYDimension; y++ {
				for x := 0; x < origXDimension; x++ {
					values[y+(i*origYDimension)][x+(j*origXDimension)] = getIncrementingTileValue(origValues[y][x], i+j)
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
			} else if x == xDimension-1 && y == yDimension-1 {
				c.goal = n
			}

			var e *graph.Edge

			if y > 0 {
				above := c.g.GetNode(geom.Pos{Y: y - 1, X: x})
				e = n.AddEdge(above, float64(above.GetProperty("value").(int)))
				e.AddProperty("dir", geom.North)
			}

			if y < yDimension-1 {
				below := c.g.GetNode(geom.Pos{Y: y + 1, X: x})
				e = n.AddEdge(below, float64(below.GetProperty("value").(int)))
				e.AddProperty("dir", geom.South)
			}

			if x > 0 {
				left := c.g.GetNode(geom.Pos{Y: y, X: x - 1})
				e = n.AddEdge(left, float64(left.GetProperty("value").(int)))
				e.AddProperty("dir", geom.West)
			}

			if x < xDimension-1 {
				right := c.g.GetNode(geom.Pos{Y: y, X: x + 1})
				e = n.AddEdge(right, float64(right.GetProperty("value").(int)))
				e.AddProperty("dir", geom.East)
			}

		}
	}

	return c
}
