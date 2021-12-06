package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reNum = regexp.MustCompile(`(\d+)`)
)

type pos struct {
	row int
	col int
}

type board struct {
	rows int
	cols int

	nums  map[int]pos
	poss  map[pos]int
	spots map[pos]bool
}

func (b *board) getEmptySum() int {
	sum := 0
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			p := pos{row: i, col: j}
			if !b.spots[p] {
				sum += b.poss[p]
			}
		}
	}
	return sum
}

func (b *board) called(num int) bool {
	if p, ok := b.nums[num]; ok {
		b.spots[p] = true

		// check all spots in this spots row
		all := true
		for i := 0; i < b.cols; i++ {
			if !b.spots[pos{row: p.row, col: i}] {
				all = false
				break
			}
		}

		if all {
			return true
		}

		// check all spots in this spots col
		all = true
		for i := 0; i < b.rows; i++ {
			if !b.spots[pos{row: i, col: p.col}] {
				all = false
				break
			}
		}

		if all {
			return true
		}
	}
	return false
}

func NewBoard(rows []string) *board {
	b := new(board)

	// num to pos map
	b.nums = make(map[int]pos)
	// spots by pos if called map
	b.spots = make(map[pos]bool)
	// pos to num map
	b.poss = make(map[pos]int)

	nums := make([][]int, len(rows), len(rows))
	for i := range rows {
		cols := reNum.FindAllStringSubmatch(rows[i], -1)
		nums[i] = make([]int, len(cols), len(cols))
		for j := range cols {
			num, _ := strconv.ParseInt(cols[j][0], 10, 32)
			nums[i][j] = int(num)
		}
	}

	b.rows = len(nums)
	b.cols = len(nums[0])

	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			p := pos{row: i, col: j}
			b.nums[nums[i][j]] = p
			b.spots[p] = false
			b.poss[p] = nums[i][j]
		}
	}

	return b
}

type game struct {
	nums   []int
	boards []*board
}

func main() {
	g := getGame()

outer:
	for _, num := range g.nums {
		for _, b := range g.boards {
			if b.called(num) {
				fmt.Println(num * b.getEmptySum())
				break outer
			}
		}
	}

}

func getGame() *game {
	lines, _ := file.GetLines("../data.txt")

	g := new(game)
	g.boards = make([]*board, 0, 100)

	nums := strings.Split(lines[0], ",")
	g.nums = make([]int, len(nums), len(nums))

	for i := range nums {
		num, _ := strconv.ParseInt(nums[i], 10, 32)
		g.nums[i] = int(num)
	}

	start := 2
	var end int

	for start < len(lines) {
		for end = start; end < len(lines) && len(lines[end]) > 0; end++ {
		}

		g.boards = append(g.boards, NewBoard(lines[start:end]))
		start = end + 1
	}

	return g
}
