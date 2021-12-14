package main

import (
	"aoc2021/common/file"
	"fmt"
	"strings"
)

type polymerTemplate struct {
	template            string
	rules               map[string]string
	production          map[string]string
	possibleProductions map[string]int
}

func newPolymerTemplate() *polymerTemplate {
	pt := new(polymerTemplate)
	pt.rules = make(map[string]string)
	pt.possibleProductions = make(map[string]int)
	pt.production = make(map[string]string)
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

func (pt *polymerTemplate) step(template string) string {
	chunkLen := len(pt.template)

	size := len(template) + len(template) - 1

	production := make([]byte, size, size)

	ptr := 0
	for i := 0; i <= len(template)-chunkLen; i += chunkLen - 1 {
		t := string(template[i : i+chunkLen])

		tprod := pt.stepChunk(t)
		copy(production[ptr:ptr+len(tprod)], []byte(tprod))

		ptr += len(tprod) - 1
	}

	return string(production)
}

func getCounts(template string) (uint64, uint64) {
	counts := make(map[byte]uint64)
	maxCount := uint64(0)
	minCount := uint64(len(template))

	for _, c := range template {
		if v, ok := counts[byte(c)]; ok {
			counts[byte(c)] = v + uint64(1)
		} else {
			counts[byte(c)] = uint64(1)
		}
	}

	for _, v := range counts {
		if v < minCount {
			minCount = v
		}
		if v > maxCount {
			maxCount = v
		}
	}

	return minCount, maxCount
}

func sub(x, y uint64) uint64 {
	return y - x
}

func main() {
	pt := getData()

	steps := 0

	template := pt.template

	chunks := make(map[string]int)
	chunksMap := make(map[int]string)
	zippedProd := make(map[int][]int)

	for {
		steps++
		fmt.Println(steps)
		prev := len(pt.production)
		template = pt.step(template)
		if prev == len(pt.production) {
			unique := 0
			for k0, v := range pt.production {
				k1 := v[0:len(pt.template)]
				k2 := v[len(pt.template)-1:]
				for _, k := range []string{k0, k1, k2} {
					if _, ok := chunks[k]; !ok {
						unique++
						chunks[k] = unique
						chunksMap[unique] = k
					}
				}
				zippedProd[chunks[k0]] = []int{chunks[k1], chunks[k2]}
			}
			break
		}
	}

	chunkLen := len(pt.template)
	zipped := make([]int, 0, len(chunks))
	for i := 0; i <= len(template)-chunkLen; i += chunkLen - 1 {
		zipped = append(zipped, chunks[template[i:i+chunkLen]])
	}

	for steps < 40 {
		steps++
		fmt.Println(steps)
		nzipped := make([]int, 0, len(zipped)+1)
		for _, i := range zipped {
			nzipped = append(nzipped, zippedProd[i]...)
		}
		zipped = nzipped
	}

	min, max := getCounts(template)

	fmt.Println(max - min)

}

func getData() *polymerTemplate {
	lines, _ := file.GetLines("../data.txt")

	pt := newPolymerTemplate()
	pt.template = lines[0]

	for _, line := range lines[2:] {
		tokens := strings.Split(line, " -> ")
		pt.rules[tokens[0]] = tokens[1]
		if v, ok := pt.possibleProductions[tokens[1]]; ok {
			pt.possibleProductions[tokens[1]] = v + 1
		} else {
			pt.possibleProductions[tokens[1]] = 1
		}
	}
	return pt
}
