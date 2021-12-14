package main

import (
	"aoc2021/common/file"
	"fmt"
	"sort"
	"strings"
)

func chunkConsecutive(input string) []string {
	chunks := make([]string, 0, len(input)-1)
	for i := 0; i < len(input)-1; i++ {
		chunk := string(input[i : i+2])
		chunks = append(chunks, chunk)
	}
	return chunks
}

func mergeProduction(input string, production []string) string {
	merged := make([]string, 0, len(input)+len(production))
	for i := 0; i < len(production); i++ {
		merged = append(merged, string(input[i]))
		merged = append(merged, production[i])
	}
	merged = append(merged, string(input[len(input)-1]))
	return strings.Join(merged, "")
}

func countElements(template string) (string, map[string]int) {
	counts := make(map[string]int)
	for _, c := range template {
		if v, ok := counts[string(c)]; ok {
			counts[string(c)] = v + 1
		} else {
			counts[string(c)] = 1
		}
	}

	chars := []byte(template)
	sort.Slice(chars, func(i, j int) bool {
		return counts[string(chars[i])] < counts[string(chars[j])]
	})
	return string(chars), counts
}

type polymerTemplate struct {
	template string
	rules    map[string]string
}

func newPolymerTemplate() *polymerTemplate {
	pt := new(polymerTemplate)
	pt.rules = make(map[string]string)
	return pt
}

func (pt *polymerTemplate) step(template string) string {
	chunks := chunkConsecutive(template)
	production := make([]string, len(chunks), len(chunks))
	for i, chunk := range chunks {
		production[i] = pt.rules[chunk]
	}
	return mergeProduction(template, production)
}

func main() {
	pt := getData()

	steps := 10

	template := pt.template

	for steps > 0 {
		template = pt.step(template)
		steps--
	}

	sorted, counts := countElements(template)

	fmt.Println(counts[string(sorted[len(sorted)-1])] - counts[string(sorted[0])])
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
