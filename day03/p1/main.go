package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

func main() {
	data := getData()
	zeros := make([]int, len(data[0]), len(data[0]))
	ones := make([]int, len(data[0]), len(data[0]))
	for _, input := range data {
		for i, c := range input {
			if c == '1' {
				ones[i]++
			} else {
				zeros[i]++
			}
		}
	}
	gamma := make([]byte, len(data[0]), len(data[0]))
	epsilon := make([]byte, len(data[0]), len(data[0]))

	for i := range ones {
		if ones[i] > zeros[i] {
			gamma[i], epsilon[i] = '1', '0'
		} else {
			gamma[i], epsilon[i] = '0', '1'
		}
	}

	g, _ := strconv.ParseInt(string(gamma), 2, 64)
	e, _ := strconv.ParseInt(string(epsilon), 2, 64)

	fmt.Println(g * e)
}

func getData() []string {
	lines, _ := file.GetLines("../data.txt")
	return lines
}
