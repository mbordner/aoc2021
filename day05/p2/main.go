package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reLine = regexp.MustCompile(`^(\d+),(\d+)\s+\->\s+(\d+),(\d+)$`)
)

type point struct {
	x int
	y int
}

func NewPoint(x, y string) point {
	p := new(point)
	i, _ := strconv.ParseInt(x, 10, 32)
	p.x = int(i)
	i, _ = strconv.ParseInt(y, 10, 32)
	p.y = int(i)
	return *p
}

type line struct {
	p1 point
	p2 point
}

func main() {
	lines := getData()

	var maxX, maxY int

	for _, line := range lines {
		if line.p1.x > maxX {
			maxX = line.p1.x
		}
		if line.p2.x > maxX {
			maxX = line.p2.x
		}
		if line.p1.y > maxY {
			maxY = line.p1.y
		}
		if line.p2.y > maxY {
			maxY = line.p2.y
		}
	}

	grid := make([][]int, maxY+1, maxY+1)
	for i := 0; i < maxY+1; i++ {
		grid[i] = make([]int, maxX+1, maxX+1)
		for j := 0; j < maxX+1; j++ {
			grid[i][j] = 0
		}
	}

	for _, l := range lines {

		var p1, p2 point
		if l.p1.x == l.p2.x || l.p1.y == l.p2.y {
			if l.p1.x <= l.p2.x && l.p1.y <= l.p2.y {
				p1 = l.p1
				p2 = l.p2
			} else {
				p1 = l.p2
				p2 = l.p1
			}
			for i := p1.y; i <= p2.y; i++ {
				for j := p1.x; j <= p2.x; j++ {
					grid[i][j]++
				}
			}
		} else {

			if l.p1.x < l.p2.x {
				p1 = l.p1
				p2 = l.p2
			} else {
				p1 = l.p2
				p2 = l.p1
			}

			dx := p2.x - p1.x
			dy := p2.y - p1.y

			m := int(float32(dy) / float32(dx))
			b := p1.y - m*p1.x

			for j := p1.x; j <= p2.x; j++ {
				i := m*j + b
				grid[i][j]++
			}

		}

	}

	count := 0
	for i := 0; i < maxY+1; i++ {
		for j := 0; j < maxX+1; j++ {
			if grid[i][j] >= 2 {
				count++
			}
		}
	}

	fmt.Println(count)
}

func getData() []line {
	lines, _ := file.GetLines("../data.txt")
	data := make([]line, 0, len(lines))
	for _, l := range lines {
		matches := reLine.FindAllStringSubmatch(l, -1)
		data = append(data, line{p1: NewPoint(matches[0][1], matches[0][2]), p2: NewPoint(matches[0][3], matches[0][4])})
	}
	return data
}
