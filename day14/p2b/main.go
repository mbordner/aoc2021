package main

import (
	"aoc2021/common/file"
	"fmt"
	"strings"
)

type polymerTemplate struct {
	template   string
	rules      map[string]string
	production map[string]string
	counts     map[string]map[byte]uint64
}

func newPolymerTemplate() *polymerTemplate {
	pt := new(polymerTemplate)
	pt.rules = make(map[string]string)
	pt.production = make(map[string]string)
	pt.counts = make(map[string]map[byte]uint64)
	return pt
}

func (pt *polymerTemplate) stepChunk(chunk string) string {
	if p, exists := pt.production[chunk]; exists {
		return p
	}

	size := len(chunk) + len(chunk) - 1

	production := make([]byte, size, size)

	for i := 0; i < len(chunk)-1; i++ {
		c1 := chunk[i]
		c3 := chunk[i+1]
		c2 := pt.rules[string([]byte{c1, c3})][0]
		production[i*2] = c1
		production[i*2+1] = c2
	}

	c3 := chunk[len(chunk)-1]
	production[len(production)-1] = c3

	pt.production[chunk] = string(production)

	return pt.production[chunk]
}

func (pt *polymerTemplate) getCounts(chunk string) map[byte]uint64 {
	counts := make(map[byte]uint64)

	if c, exists := pt.counts[chunk]; exists {
		for k, v := range c {
			counts[k] = v
		}
	} else {
		c := make(map[byte]uint64)

		for _, b := range chunk {
			if v, ok := c[byte(b)]; ok {
				c[byte(b)] = v + uint64(1)
			} else {
				c[byte(b)] = uint64(1)
			}
		}

		pt.counts[chunk] = c

		for k, v := range c {
			counts[k] = v
		}
	}

	return counts
}

func (pt *polymerTemplate) step(template string, steps int) map[byte]uint64 {
	steps--

	if steps >= 0 {
		chunkLen := len(pt.template)

		production := pt.stepChunk(template)

		leftCounts := pt.step(production[0:chunkLen], steps)
		rightCounts := pt.step(production[chunkLen-1:], steps)

		leftCounts[production[chunkLen-1]]--

		for k, v := range rightCounts {
			leftCounts[k] += v
		}

		return leftCounts
	} else {
		return pt.getCounts(template)
	}

}

func main() {
	pt := getData()

	counts := pt.step(pt.template, 40)

	for k, v := range counts {
		fmt.Println(string(k), ":", v)
	}

	var max, min uint64
	for _, v := range counts {
		min += v
		if v > max {
			max = v
		}
	}
	for _, v := range counts {
		if v < min {
			min = v
		}
	}
	fmt.Println(max - min)
}

func getData() *polymerTemplate {
	lines, _ := file.GetLines("../data.txt")

	pt := newPolymerTemplate()
	pt.template = lines[0]

	for _, line := range lines[2:] {
		tokens := strings.Split(line, " -> ")
		pt.rules[tokens[0]] = tokens[1]
	}
	return pt
}
