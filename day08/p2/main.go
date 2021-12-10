package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reEntry = regexp.MustCompile(`^([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)\s\|\s([a-g]+)\s([a-g]+)\s([a-g]+)\s([a-g]+)$`)

	vals = map[string]int{
		"a": 0b0000001,
		"b": 0b0000010,
		"c": 0b0000100,
		"d": 0b0001000,
		"e": 0b0010000,
		"f": 0b0100000,
		"g": 0b1000000,
	}

	valZero  = vals["a"] | vals["b"] | vals["c"] | vals["e"] | vals["f"] | vals["g"]
	valOne   = vals["c"] | vals["f"]
	valTwo   = vals["a"] | vals["c"] | vals["d"] | vals["e"] | vals["g"]
	valThree = vals["a"] | vals["c"] | vals["d"] | vals["f"] | vals["g"]
	valFour  = vals["b"] | vals["c"] | vals["d"] | vals["f"]
	valFive  = vals["a"] | vals["b"] | vals["d"] | vals["f"] | vals["g"]
	valSix   = vals["a"] | vals["b"] | vals["d"] | vals["e"] | vals["f"] | vals["g"]
	valSeven = vals["a"] | vals["c"] | vals["f"]
	valEight = vals["a"] | vals["b"] | vals["c"] | vals["d"] | vals["e"] | vals["f"] | vals["g"]
	valNine  = vals["a"] | vals["b"] | vals["c"] | vals["d"] | vals["f"] | vals["g"]
)

func intersect(x, y string) (xo, common, yo string) {
	cmap := make(map[byte]int)
	for _, c := range []byte(x) {
		cmap[c] = 1
	}
	for _, c := range []byte(y) {
		if _, ok := cmap[c]; ok {
			cmap[c]++
		} else {
			yo += string(c)
		}
	}
	for b, c := range cmap {
		if c == 1 {
			xo += string(b)
		} else {
			common += string(b)
		}
	}
	return
}

type entry struct {
	patterns []string
	digits   []string
	mapped   map[string]string
	value    int
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
	entr := new(entry)
	matches := reEntry.FindAllStringSubmatch(line, -1)
	entr.patterns = matches[0][1:11]
	entr.digits = matches[0][11:]

	entr.mapped = make(map[string]string)

	var one, seven, four, eight string

	for _, p := range entr.patterns {
		switch len(p) {
		case 2:
			one = p
		case 3:
			seven = p
		case 4:
			four = p
		case 7:
			eight = p
		}
	}

	_, tmp, a := intersect(one, seven)
	entr.mapped[a] = "a"

	_, _, bd := intersect(tmp, four)

	_, _, eg := intersect(seven+four, eight)

	var zero, six, nine string

	for _, p := range entr.patterns {
		if len(p) == 6 {
			_, eORg, _ := intersect(bd, p)
			_, bORd, _ := intersect(eg, p)
			if len(eORg) == 2 && len(bORd) == 2 {
				six = p
			} else if len(eORg) == 2 {
				nine = p
			} else if len(bORd) == 2 {
				zero = p
			}
		}
	}

	c, _, _ := intersect(seven, six)

	entr.mapped[c] = "c"

	f, _, _ := intersect(one, c)

	entr.mapped[f] = "f"

	d, _, _ := intersect(eight, zero)

	entr.mapped[d] = "d"

	var two, three, five string

	for _, p := range entr.patterns {
		if len(p) == 5 {
			_, _, tmp = intersect(p, nine)

			if tmp == c {
				five = p
			} else if len(tmp) == 2 {
				two = p
			} else {
				three = p
			}

		}
	}

	b, _, _ := intersect(eight, two+f)

	entr.mapped[b] = "b"

	e, _, _ := intersect(eight, five+c)

	entr.mapped[e] = "e"

	used := ""

	for k := range entr.mapped {
		used += k
	}

	g, _, _ := intersect(eight, used)

	entr.mapped[g] = "g"

	check := ""
	for _, b := range three {
		check += entr.mapped[string(b)]
	}
	tmp, _, _ = intersect("acdfg", check)

	if len(tmp) != 0 {
		panic("not right")
	}

	valStr := ""
	for _, digit := range entr.digits {
		val := 0
		for _, b := range digit {
			val |= vals[entr.mapped[string(b)]]
		}
		switch val {
		case valZero:
			valStr += "0"
		case valOne:
			valStr += "1"
		case valTwo:
			valStr += "2"
		case valThree:
			valStr += "3"
		case valFour:
			valStr += "4"
		case valFive:
			valStr += "5"
		case valSix:
			valStr += "6"
		case valSeven:
			valStr += "7"
		case valEight:
			valStr += "8"
		case valNine:
			valStr += "9"
		}
	}

	tmpVal, _ := strconv.ParseInt(valStr, 10, 32)
	entr.value = int(tmpVal)

	return entr
}

func main() {
	entries := getData()
	count := 0
	for _, e := range entries {
		count += e.value
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
