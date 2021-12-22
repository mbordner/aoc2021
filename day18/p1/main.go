package main

import (
	"aoc2021/common/file"
	"fmt"
	"regexp"
	"strconv"
)

var (
	reDigit = regexp.MustCompile(`^(\d+)$`)
)

type sfNum struct {
	parent *sfNum
	value  int
	left   *sfNum
	right  *sfNum
}

func (sfn *sfNum) String() string {
	if sfn.isLeaf() {
		return fmt.Sprintf("%d", sfn.getValue())
	}
	return fmt.Sprintf("[%s,%s]", sfn.left, sfn.right)
}

func (sfn *sfNum) add(osfn *sfNum) *sfNum {
	nsfn := &sfNum{}
	sfn.parent = nsfn
	osfn.parent = nsfn
	nsfn.left = sfn
	nsfn.right = osfn
	return nsfn
}

func (sfn *sfNum) getLeafs() []*sfNum {
	sfns := make([]*sfNum, 0, 10)
	if sfn.isLeaf() {
		sfns = append(sfns, sfn)
	} else {
		sfns = append(sfns, sfn.left.getLeafs()...)
		sfns = append(sfns, sfn.right.getLeafs()...)
	}
	return sfns
}

func (sfn *sfNum) getRoot() *sfNum {
	p := sfn
	for p.parent != nil {
		p = p.parent
	}
	return p
}

func (sfn *sfNum) reduce() {

	toE := sfn.toExplode(1)
	if len(toE) > 0 {

		pair := toE[0]
		isLeft := pair.parent.left == pair
		root := sfn.getRoot()
		leafs := root.getLeafs()
		var index int
		for index = 0; index < len(leafs); index++ {
			if pair.left == leafs[index] {
				break
			}
		}
		if index > 0 {
			leafs[index-1].addValue(pair.left.getValue())
		}
		if index < len(leafs)-2 {
			leafs[index+2].addValue(pair.right.getValue())
		}
		nsfn := &sfNum{}
		nsfn.value = 0
		nsfn.parent = pair.parent
		if isLeft {
			pair.parent.left = nsfn
		} else {
			pair.parent.right = nsfn
		}
		pair.parent = nil

		sfn.reduce()
	}

	toS := sfn.toSplit()
	if len(toS) > 0 {
		leaf := toS[0]
		isLeft := leaf.parent.left == leaf

		leftNum := leaf.value / 2
		rightNum := leaf.value/2 + leaf.value%2

		nsfn := newSFNum(fmt.Sprintf("[%d,%d]", leftNum, rightNum))
		nsfn.parent = leaf.parent

		if isLeft {
			leaf.parent.left = nsfn
		} else {
			leaf.parent.right = nsfn
		}

		leaf.parent = nil

		sfn.reduce()
	}

}

func (sfn *sfNum) addAndReduce(osfn *sfNum) *sfNum {
	nsfn := sfn.add(osfn)
	nsfn.reduce()
	return nsfn
}

func (sfn *sfNum) toExplode(depth int) []*sfNum {
	sfns := make([]*sfNum, 0, 10)
	if depth > 4 {
		if sfn.isLeaf() == false {
			sfns = append(sfns, sfn)
		}
	} else if sfn.isLeaf() == false {
		sfns = append(sfns, sfn.left.toExplode(depth+1)...)
		sfns = append(sfns, sfn.right.toExplode(depth+1)...)
	}
	return sfns
}

func (sfn *sfNum) toSplit() []*sfNum {
	sfns := make([]*sfNum, 0, 10)
	for _, tsfn := range []*sfNum{sfn.left, sfn.right} {
		if tsfn.isLeaf() {
			if tsfn.getValue() >= 10 {
				sfns = append(sfns, tsfn)
			}
		} else {
			sfns = append(sfns, tsfn.toSplit()...)
		}
	}
	return sfns
}

func (sfn *sfNum) getMagnitude() int {
	if sfn.isLeaf() {
		return sfn.getValue()
	}
	return 3*sfn.left.getMagnitude() + 2*sfn.right.getMagnitude()
}

func (sfn *sfNum) isLeaf() bool {
	if sfn.left != nil && sfn.right != nil {
		return false
	}
	return true
}

func (sfn *sfNum) addValue(value int) {
	sfn.value += value
}

func (sfn *sfNum) getValue() int {
	return sfn.value
}

func (sfn *sfNum) getDepth() int {
	depth := 1
	ld := 0
	rd := 0
	if sfn.left.isLeaf() == false {
		ld = sfn.left.getDepth()
	}
	if sfn.right.isLeaf() == false {
		rd = sfn.right.getDepth()
	}
	if ld > rd {
		depth += ld
	} else {
		depth += rd
	}
	return depth
}

func newSFNum(s string) *sfNum {
	sfn := new(sfNum)

	if reDigit.MatchString(s) {
		val, _ := strconv.ParseInt(s, 10, 32)
		sfn.value = int(val)
	} else {
		if s[0] != '[' && s[len(s)-1] != ']' {
			panic("syntax error : expected pair to be wrapped in [ ]")
		}
		s := s[1 : len(s)-1]

		start := 0
		end := 0
		for side := 0; side < 2; side++ {
			if s[start] == '[' {
				stack := 1
				for end = start + 1; stack > 0; end++ {
					if s[end] == '[' {
						stack++
					} else if s[end] == ']' {
						stack--
					}
				}
			} else {
				for end = start; end < len(s) && reDigit.MatchString(string(s[end])); end++ {
				}
			}
			sideSFNum := newSFNum(s[start:end])
			sideSFNum.parent = sfn
			if side == 0 {
				sfn.left = sideSFNum
				if s[end] != ',' {
					panic("syntax error parsing left side of pair")
				}
				start = end + 1
			} else {
				sfn.right = sideSFNum
				if end != len(s) {
					panic("syntax error parsing right side of pair")
				}
			}
		}

	}

	return sfn
}

func getLargestMagnitudeAddingOnlyTwo(lines []string) int {
	max := 0
	var sfn *sfNum
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines); j++ {
			if j != i {
				sfn = newSFNum(lines[i])
				sfn = sfn.addAndReduce(newSFNum(lines[j]))
				mag := sfn.getMagnitude()
				if mag > max {
					max = mag
				}
			}
		}
	}
	return max
}

func getList(nums []string) []*sfNum {
	list := make([]*sfNum, 0, len(nums))
	for _, num := range nums {
		list = append(list, newSFNum(num))
	}
	return list
}

func getData(filename string) []string {
	lines, _ := file.GetLines(filename)
	return lines
}

func main() {
	lines := getData("../data.txt")
	list := getList(lines)
	sfn := list[0]
	for _, nsfn := range list[1:] {
		sfn = sfn.addAndReduce(nsfn)
	}
	fmt.Println(sfn.getMagnitude())

	max := getLargestMagnitudeAddingOnlyTwo(lines)

	fmt.Println(max)
}
