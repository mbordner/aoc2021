package main

import (
	"aoc2021/common/file"
	"aoc2021/day16/packet"
	"fmt"
)

func main() {
	transmission := getTransmission("../data.txt")

	p := packet.Decompile(transmission)

	fmt.Println(p)

	fmt.Println(p.Eval())
}

func getTransmission(filename string) string {
	lines, _ := file.GetLines(filename)
	return lines[0]
}
