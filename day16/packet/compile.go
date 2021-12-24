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

func createLiteralBitString(value uint64) string {
	bitsPerGroup := 4
	groupSize := bitsPerGroup + 1

	bits := reverse(fmt.Sprintf("%b", value))
	l := len(bits)
	pl := l / bitsPerGroup * groupSize
	if l%bitsPerGroup > 0 {
		pl += groupSize
	}

	// create bits , and initialize them to 0
	payload := make([]byte, pl, pl)
	for i := 0; i < len(payload); i++ {
		payload[i] += byte('0')
	}

	marker := byte('0')
	for i, j := 0, 0; i < l; i, j = i+1, j+1 {
		payload[j] = bits[i]
		if i != 0 && i%bitsPerGroup == 0 {
			j++
			payload[j] = marker
			marker = byte('1')
		}
	}

	payload = []byte(reverse(string(payload)))
	payload[0] = '1'

	return string(payload)
}

func getSubpacketSourceStrings(src string) ([]string, error) {
	srcs := make([]string, 0, 10)

	start := 0
	end := 0

	for start < len(src) {

		if reDigit.MatchString(string(src[end])) {
			for end = start; end < len(src); end++ {
				if !reDigit.MatchString(string(src[end])) {
					break
				}
			}
		} else if src[start] != '(' {
			return nil, errors.New("expected (")
		} else {
			end = start
			stack := 1
			for stack > 0 && end < len(src) {
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
			end++
		}

		srcs = append(srcs, src[start:end])

		start = end
		if start < len(src) {
			if !reWhiteSpace.MatchString(string(src[start])) {
				return nil, errors.New("expected whitespace between sub packets")
			}
			start++
		}

	}

	return srcs, nil
}

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
		for stack > 0 && end < len(src) {
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
		p.payload = createLiteralBitString(value)
	} else {
		p.header = newPacketHeaderFromValues(typeID, 1)

		subpacketSourceStrings, err := getSubpacketSourceStrings(src[end+1:])
		if err != nil {
			return nil, err
		}

		p.subpackets = make([]*Packet, 0, len(subpacketSourceStrings))

		spLen := 0

		for _, spsrc := range subpacketSourceStrings {
			sp, err := parseTree(spsrc)
			if err != nil {
				return nil, err
			}
			p.subpackets = append(p.subpackets, sp)
			spLen += sp.length()
		}

		if spLen <= 32768 {
			p.payload = "0" + fmt.Sprintf("%015b", spLen)
		} else {
			p.payload = "1" + fmt.Sprintf("%011b", len(p.subpackets))
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
	for i, j := 0, len(b)-1; i <= lh; i, j = i+1, j-1 {
		b[i], b[j] = s[j], s[i]
	}
	return string(b)
}
