package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

var (
	maxEnergy = 9
)

type pos struct {
	y int
	x int
}

func getAdjacents(grid [][]int, p pos) []pos {
	adj := make([]pos, 0, 8)
	if p.x > 0 { // left
		adj = append(adj, pos{y: p.y, x: p.x - 1})

		if p.y > 0 { // upper left
			adj = append(adj, pos{y: p.y - 1, x: p.x - 1})
		}
		if p.y < len(grid)-1 { // lower left
			adj = append(adj, pos{y: p.y + 1, x: p.x - 1})
		}
	}
	if p.y > 0 { // above
		adj = append(adj, pos{y: p.y - 1, x: p.x})
	}
	if p.x < len(grid[p.y])-1 { // right
		adj = append(adj, pos{y: p.y, x: p.x + 1})

		if p.y > 0 { // upper right
			adj = append(adj, pos{y: p.y - 1, x: p.x + 1})
		}
		if p.y < len(grid)-1 { // lower right
			adj = append(adj, pos{y: p.y + 1, x: p.x + 1})
		}
	}
	if p.y < len(grid)-1 { // below
		adj = append(adj, pos{y: p.y + 1, x: p.x})
	}
	return adj
}

func flash(ready map[pos]int, grid [][]int) map[pos]int {
	newReady := make(map[pos]int)

	for p := range ready {
		adjs := getAdjacents(grid, p)

		for _, a := range adjs {
			grid[a.y][a.x]++
			if grid[a.y][a.x] == maxEnergy+1 {
				newReady[a] = grid[a.y][a.x]
			}
		}
	}

	return newReady
}

func main() {
	grid := getData()

	flashes := 0

	rounds := 0

	for true {

		rounds++

		max := make(map[pos]int)

		for i := 0; i < len(grid); i++ {
			for j := 0; j < len(grid[i]); j++ {
				grid[i][j]++
				if grid[i][j] == maxEnergy+1 {
					max[pos{y: i, x: j}] = grid[i][j]
				}
			}
		}

		if len(max) > 0 {

			newMax := flash(max, grid)

			for len(newMax) > 0 {
				for p, v := range newMax {
					max[p] = v
				}
				newMax = flash(newMax, grid)

			}

			flashes += len(max)

			for p := range max {
				grid[p.y][p.x] = 0
			}

			if len(max) == len(grid)*len(grid[0]) {
				fmt.Println("all flash at ", rounds)
				break
			}
		}

	}

}

func getData() [][]int {
	lines, _ := file.GetLines("../data.txt")
	data := make([][]int, len(lines), len(lines))
	for i, line := range lines {
		data[i] = make([]int, len(line), len(line))
		for j, d := range line {
			v, _ := strconv.ParseInt(string(d), 10, 32)
			data[i][j] = int(v)
		}
	}
	return data
}
