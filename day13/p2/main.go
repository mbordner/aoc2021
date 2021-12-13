package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
	"strings"
)

type pos struct {
	y int
	x int
}

func (p pos) String() string {
	return fmt.Sprintf("(%d,%d)", p.x, p.y)
}

type fold struct {
	axis  string
	value int
}

type sheet struct {
	points           map[pos]bool
	foldInstructions []*fold
	maxY             int
	maxX             int
}

func newSheet() *sheet {
	s := new(sheet)
	s.points = make(map[pos]bool)
	s.foldInstructions = make([]*fold, 0, 100)
	return s
}

func (s *sheet) add(p pos) {
	if p.x > s.maxX {
		s.maxX = p.x
	}
	if p.y > s.maxY {
		s.maxY = p.y
	}
	s.points[p] = true
}

func (s *sheet) hasPoint(p pos) bool {
	if b, exists := s.points[p]; exists {
		return b
	}
	return false
}

func (s *sheet) getGrid() [][]string {
	grid := make([][]string, s.maxY+1, s.maxY+1)
	for i := 0; i <= s.maxY; i++ {
		grid[i] = make([]string, s.maxX+1, s.maxX+1)
		for j := 0; j <= s.maxX; j++ {
			if s.hasPoint(pos{y: i, x: j}) {
				grid[i][j] = "#"
			} else {
				grid[i][j] = "."
			}
		}
	}
	return grid
}

func (s *sheet) printGrid(g [][]string) {
	for _, r := range g {
		fmt.Println(strings.Join(r, ""))
	}
}

func (s *sheet) getCount() int {
	return len(s.points)
}

func (s *sheet) fold(a string, v int) {
	grid := s.getGrid()
	n := newSheet()
	if a == "y" {
		for i := 1; i+v < len(grid) && v-i >= 0; i++ {
			for j := 0; j < len(grid[i]); j++ {
				if grid[v+i][j] == "#" {
					grid[v-i][j] = "#"
				}
			}
		}
		for i := 0; i < v; i++ {
			for j := 0; j < len(grid[i]); j++ {
				if grid[i][j] == "#" {
					n.add(pos{y: i, x: j})
				}
			}
		}
	} else {
		for j := 1; v+j < len(grid[0]) && v-j >= 0; j++ {
			for i := 0; i < len(grid); i++ {
				if grid[i][v+j] == "#" {
					grid[i][v-j] = "#"
				}
			}
		}
		for j := 0; j < v; j++ {
			for i := 0; i < len(grid); i++ {
				if grid[i][j] == "#" {
					n.add(pos{y: i, x: j})
				}
			}
		}
	}

	s.points = n.points
	s.maxY = n.maxY
	s.maxX = n.maxX
}

func main() {
	s := getSheet("../data.txt")

	for _, f := range s.foldInstructions {
		s.fold(f.axis, f.value)
	}

	s.printGrid(s.getGrid())
}

func getInt(i string) int {
	v, _ := strconv.ParseInt(i, 10, 32)
	return int(v)
}

func getSheet(filename string) *sheet {
	lines, _ := file.GetLines(filename)

	s := newSheet()

	donePoints := false
	for _, line := range lines {
		if line == "" {
			donePoints = true
			continue
		}
		if donePoints {
			f := new(fold)

			tokens := strings.Split(line, "=")
			f.value = getInt(tokens[1])

			tokens = strings.Split(tokens[0], " ")
			f.axis = tokens[2]

			s.foldInstructions = append(s.foldInstructions, f)
		} else {

			tokens := strings.Split(line, ",")

			p := new(pos)
			p.x = getInt(tokens[0])
			p.y = getInt(tokens[1])

			s.add(*p)
		}
	}

	return s
}
