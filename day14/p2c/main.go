package main

import (
	"aoc2021/common/file"
	"fmt"
	"strings"
)

type freqCounts map[string]uint64

type polymerTemplate struct {
	template string
	rules    map[string]string
}

func newPolymerTemplate() *polymerTemplate {
	pt := new(polymerTemplate)
	pt.rules = make(map[string]string)
	return pt
}

func (pt *polymerTemplate) step(counts freqCounts) freqCounts {
	newCounts := make(freqCounts)
	for pair, count := range counts {
		production := pt.rules[pair]
		pair1 := string(pair[0]) + production
		pair2 := production + string(pair[1])
		for _, p := range []string{pair1, pair2} {
			if v, exists := newCounts[p]; exists {
				newCounts[p] = v + count
			} else {
				newCounts[p] = count
			}
		}
	}
	return newCounts
}

func (pt *polymerTemplate) getCounts(pairCounts freqCounts) map[byte]uint64 {
	counts := make(map[byte]uint64)

	for k, v := range pairCounts {
		b := k[0]
		if c, exists := counts[b]; exists {
			counts[b] = c + v
		} else {
			counts[b] = v
		}
	}

	lastc := pt.template[len(pt.template)-1]

	if c, exists := counts[lastc]; exists {
		counts[lastc] = c + 1
	} else {
		counts[lastc] = uint64(1)
	}

	return counts
}

func main() {
	pt := getData()

	freqs := make(freqCounts)

	l := len(pt.template)
	for i := 0; i < l-1; i++ {
		combn := pt.template[i : i+2]
		if v, exists := freqs[combn]; exists {
			freqs[combn] = v + uint64(1)
		} else {
			freqs[combn] = uint64(1)
		}
	}

	steps := 40

	for steps > 0 {
		steps--
		freqs = pt.step(freqs)
	}

	counts := pt.getCounts(freqs)

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
