package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
)

var (
	reEntry = regexp.MustCompile(`^([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s\|\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)$`)
)

type entry struct {
	patterns []string
	digits   []string
}

func (e *entry) uniqueCount() int {
	count := 0
	for _, v := range e.digits {
		if len(v) == 2 || len(v) == 3 || len(v) == 4 || len(v) == 7 {
			count++
		}
	}
	return count
}

func NewEntry(line string) *entry {
	e := new(entry)
	matches := reEntry.FindAllStringSubmatch(line, -1)
	e.patterns = matches[0][1:11]
	e.digits = matches[0][11:]
	return e
}

func main() {
	entries := getData()
	count := 0
	for _, e := range entries {
		count += e.uniqueCount()
	}
	fmt.Println(count)
}

func getData() []*entry {
	lines, _ := file.GetLines("../data.txt")
	entries := make([]*entry, 0, len(lines))
	for _, line := range lines {
		e := NewEntry(line)
		entries = append(entries, e)
	}
	return entries
}

/**

  0:      1:      2:      3:      4:
 aaaa    ....    aaaa    aaaa    ....
b    c  .    c  .    c  .    c  b    c
b    c  .    c  .    c  .    c  b    c
 ....    ....    dddd    dddd    dddd
e    f  .    f  e    .  .    f  .    f
e    f  .    f  e    .  .    f  .    f
 gggg    ....    gggg    gggg    ....

  5:      6:      7:      8:      9:
 aaaa    aaaa    aaaa    aaaa    aaaa
b    .  b    .  .    c  b    c  b    c
b    .  b    .  .    c  b    c  b    c
 dddd    dddd    ....    dddd    dddd
.    f  e    f  .    f  e    f  .    f
.    f  e    f  .    f  e    f  .    f
 gggg    gggg    ....    gggg    gggg


*/
