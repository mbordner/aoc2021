package main

import (
	"aoc2021/common/file"
	"aoc2021/day16/packet"
	"fmt"
)

func main() {
	var p *packet.Packet

	tests := []struct {
		input  string
		output uint64
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}

	for _, test := range tests {
		p = packet.Decompile(test.input)
		fmt.Println(p)
		if p.Eval() != test.output {
			panic("doesn't pass")
		}
	}

	transmission := getTransmission("../data.txt")

	p = packet.Decompile(transmission)

	fmt.Println(p)

	fmt.Println(p.Eval())

}

func getTransmission(filename string) string {
	lines, _ := file.GetLines(filename)
	return lines[0]
}
