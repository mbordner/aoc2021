package main

import (
	"aoc2021/common/file"
	"fmt"
	"sort"
	"strconv"
)

type pos struct {
	y int
	x int
}

type basin map[pos]int

func getBasin(points [][]int, p pos) basin {
	b := make(basin)
	adj := getAdjacents(points, p)
	higher := make([]pos, 0, len(adj))
	for _, a := range adj {
		if points[a.y][a.x] > points[p.y][p.x] && points[a.y][a.x] != 9 {
			higher = append(higher, a)
		}
	}
	b[p] = points[p.y][p.x]
	if len(higher) > 0 {
		for _, h := range higher {
			hb := getBasin(points, h)
			for k, v := range hb {
				b[k] = v
			}
		}
	}
	return b
}

func getAdjacents(points [][]int, p pos) []pos {
	adj := make([]pos, 0, 4)
	if p.x > 0 { // left
		adj = append(adj, pos{y: p.y, x: p.x - 1})
	}
	if p.y > 0 { // above
		adj = append(adj, pos{y: p.y - 1, x: p.x})
	}
	if p.x < len(points[p.y])-1 { // right
		adj = append(adj, pos{y: p.y, x: p.x + 1})
	}
	if p.y < len(points)-1 { // below
		adj = append(adj, pos{y: p.y + 1, x: p.x})
	}
	return adj
}

func main() {
	points := getData()
	riskSum := 0

	lowPoints := make([]pos, 0, 100)

	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points[i]); j++ {
			p := points[i][j]
			adj := getAdjacents(points, pos{y: i, x: j})
			lowPoint := true
			for _, a := range adj {
				if p >= points[a.y][a.x] {
					lowPoint = false
					break
				}
			}
			if lowPoint {
				risk := 1 + p
				riskSum += risk
				lowPoints = append(lowPoints, pos{y: i, x: j})
			}
		}
	}
	fmt.Println(riskSum)

	basins := make([]basin, 0, len(lowPoints))

	for _, lp := range lowPoints {
		b := getBasin(points, lp)
		basins = append(basins, b)
		//fmt.Println(lp, len(b), b)
	}

	sort.Slice(basins, func(i, j int) bool {
		return len(basins[i]) > len(basins[j])
	})

	value := len(basins[0])
	for _, b := range basins[1:3] {
		value *= len(b)
	}
	fmt.Println(value)

}

func getData() [][]int {
	lines, _ := file.GetLines("../data.txt")
	points := make([][]int, 0, len(lines))
	for _, line := range lines {
		ps := make([]int, 0, len(line))
		for _, p := range line {
			h, _ := strconv.Atoi(string(p))
			ps = append(ps, h)
		}
		points = append(points, ps)
	}
	return points
}
