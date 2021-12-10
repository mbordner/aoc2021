package main

import (
	"aoc2021/common/file"
	"fmt"
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

	for _, line := range lines {
		s := newStack()
		for _, c := range line {
			if isStart(string(c)) {
				s.push(string(c))
			} else {
				if end[string(c)] == s.peek() {
					s.pop()
				} else {
					errors[string(c)]++
					break
				}
			}
		}
	}

	score := 0

	for c, v := range errors {
		if v > 0 {
			switch c {
			case ")":
				score += v * 3
			case "]":
				score += v * 57
			case "}":
				score += v * 1197
			case ">":
				score += v * 25137
			}
		}
	}

	fmt.Println(score)

}

func getData() []string {
	lines, _ := file.GetLines("../data.txt")
	return lines
}
