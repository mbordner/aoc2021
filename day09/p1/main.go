package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

func main() {
	points := getData()
	riskSum := 0
	for i := 0; i < len(points); i++ {
		for j := 0; j < len(points[i]); j++ {
			p := points[i][j]
			adj := make([]int, 0, 4)
			if j > 0 { // left
				adj = append(adj, points[i][j-1])
			}
			if i > 0 { // above
				adj = append(adj, points[i-1][j])
			}
			if j < len(points[i])-1 { // right
				adj = append(adj, points[i][j+1])
			}
			if i < len(points) - 1 { // below
				adj = append(adj, points[i+1][j])
			}
			lowPoint := true
			for _, a := range adj {
				if p >= a {
					lowPoint = false
					break
				}
			}
			if lowPoint {
				risk := 1 + p
				riskSum += risk
			}
		}
	}
	fmt.Println(riskSum)
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
