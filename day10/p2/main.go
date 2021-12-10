package main

import (
	"aoc2021/common/file"
	"fmt"
	"sort"
)

var (
	start = map[string]string{
		"{": "}",
		"[": "]",
		"(": ")",
		"<": ">",
	}
	end = map[string]string{
		"}": "{",
		"]": "[",
		")": "(",
		">": "<",
	}
)

func isStart(s string) bool {
	if _, ok := start[s]; ok {
		return true
	}
	return false
}

func isEnd(s string) bool {
	if _, ok := end[s]; ok {
		return true
	}
	return false
}

type stack struct {
	s []string
}

func (s *stack) push(c string) {
	s.s = append(s.s, c)
}

func (s *stack) peek() string {
	if len(s.s) > 0 {
		return s.s[len(s.s)-1]
	}
	return ""
}

func (s *stack) pop() string {
	c := s.peek()
	if c != "" {
		s.s = s.s[0 : len(s.s)-1]
	}
	return c
}

func newStack() *stack {
	s := new(stack)
	s.s = make([]string, 0, 100)
	return s
}

func main() {
	lines := getData()

	errors := map[string]int{
		"}": 0,
		")": 0,
		"]": 0,
		">": 0,
	}

	incomplete := make([]string, 0, len(lines))

	for _, line := range lines {
		s := newStack()
		errored := false
		for _, c := range line {
			if isStart(string(c)) {
				s.push(string(c))
			} else {
				if end[string(c)] == s.peek() {
					s.pop()
				} else {
					errors[string(c)]++
					errored = true
					break
				}
			}
		}
		if !errored && s.peek() != "" {
			incomplete = append(incomplete, line)
		}
	}

	syntaxErrorScore := 0

	for c, v := range errors {
		if v > 0 {
			switch c {
			case ")":
				syntaxErrorScore += v * 3
			case "]":
				syntaxErrorScore += v * 57
			case "}":
				syntaxErrorScore += v * 1197
			case ">":
				syntaxErrorScore += v * 25137
			}
		}
	}

	fmt.Println(syntaxErrorScore)

	completes := make([]string, 0, len(incomplete))

	autoCompleteScores := make([]int, 0, len(incomplete))

	for _, line := range incomplete {
		autoCompleteScore := 0
		s := newStack()
		for _, c := range line {
			if isStart(string(c)) {
				s.push(string(c))
			} else {
				if end[string(c)] == s.peek() {
					s.pop()
				} else {
					panic("not supposed to happen")
				}
			}
		}
		complete := ""
		for s.peek() != "" {
			c := s.pop()
			autoCompleteScore *= 5
			switch start[c] {
			case ")":
				autoCompleteScore += 1
			case "]":
				autoCompleteScore += 2
			case "}":
				autoCompleteScore += 3
			case ">":
				autoCompleteScore += 4
			}
			complete += start[c]
		}

		completes = append(completes, complete)
		autoCompleteScores = append(autoCompleteScores, autoCompleteScore)
	}

	sort.Slice(autoCompleteScores, func(i, j int) bool {
		return autoCompleteScores[i] < autoCompleteScores[j]
	})

	fmt.Println(autoCompleteScores[len(autoCompleteScores)/2])

}

func getData() []string {
	lines, _ := file.GetLines("../data.txt")
	return lines
}
