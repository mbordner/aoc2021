package packet

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reWhiteSpace = regexp.MustCompile(`\s+`)
	reDigit      = regexp.MustCompile(`\d+`)
)

func parseTree(src string) (*Packet, error) {
	var typeID Operation
	p := &Packet{}

	start := 0
	end := 0

	if reDigit.MatchString(string(src[end])) {
		for end = start; end < len(src); end++ {
			if !reDigit.MatchString(string(src[end])) {
				break
			}
		}
		typeID = Literal
	} else if src[start] != '(' {
		return nil, errors.New("expected (")
	} else {
		end = start
		stack := 1
		for stack == 0 && end < len(src) {
			end++
			if src[end] == '(' {
				stack++
			} else if src[end] == ')' {
				stack--
			}
		}
		if src[end] != ')' {
			return nil, errors.New("exepcted )")
		}
		src = src[start+1 : end]

		start = 0
		end = start
		for end < len(src) && !reWhiteSpace.MatchString(string(src[end])) {
			end++
		}

		if reDigit.MatchString(src[start:end]) {
			typeID = Literal
		} else {
			op, err := OperationFromString(src[start:end])
			if err != nil {
				return nil, err
			}
			typeID = op

			if end == len(src) {
				return nil, errors.New("missing operation params")
			}
		}
	}

	if typeID == Literal {
		p.header = newPacketHeaderFromValues(Literal, 1)
		value, err := strconv.ParseUint(src[start:end], 10, 64)
		if err != nil {
			return nil, err
		}
		bits := reverse(fmt.Sprintf("%b", value))
		payload := make([]byte, 0, 4*(len(bits)/3+len(bits)%3))
		for i, j := 0, len(bits)-1; j >= 0; i, j = i+1, j-1 {
			if i%3 == 0 {

			}
		}
		if payload != nil {

		}
	}

	return p, nil
}

func Compile(src string) (*Packet, error) {
	src = reWhiteSpace.ReplaceAllString(src, " ")
	src = strings.TrimSpace(src)

	return parseTree(src)
}

func reverse(s string) string {
	b := make([]byte, len(s), len(s))
	lh := len(s) / 2
	for i, j := 0, len(b)-1; i < lh; i, j = i+1, j-1 {
		b[i], b[j] = s[j], s[i]
	}
	return string(b)
}
