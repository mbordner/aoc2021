package main

import (
	"aoc2021/common/file"
	"errors"
	"fmt"
	"strconv"
)

func main() {
	var packets []*packet

	tests := []struct {
		input  string
		output uint64
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}

	for _, test := range tests {
		packets = parse(binary(test.input), -1)
		if packets[0].eval() != test.output {
			panic("doesn't pass")
		}
	}

	transmission := getTransmission("../data.txt")

	packets = parse(transmission, -1)

	fmt.Println(packets[0].eval())

}

type packetHeader struct {
	version byte
	typeID  byte
	payload string
}

func newPacketHeader(value string) *packetHeader {
	ph := new(packetHeader)
	ph.payload = value

	val, _ := strconv.ParseInt(value[0:3], 2, 8)
	ph.version = byte(val)

	val, _ = strconv.ParseInt(value[3:], 2, 8)
	ph.typeID = byte(val)

	return ph
}

type packet struct {
	header     *packetHeader
	payload    string
	subpackets []*packet
}

func (p *packet) length() int {
	l := len(p.header.payload) + len(p.payload)
	if p.subpackets != nil {
		for _, sp := range p.subpackets {
			l += sp.length()
		}
	}
	return l
}

func (p *packet) literal() uint64 {
	value := make([]byte, 0, len(p.payload))
	for i := 0; i < len(p.payload); i += 5 {
		value = append(value, p.payload[i+1:i+5]...)
	}
	val, _ := strconv.ParseUint(string(value), 2, 64)
	return val
}

func (p *packet) sum() uint64 {
	val := uint64(0)
	for _, sp := range p.subpackets {
		spval := sp.eval()
		val += spval
	}
	return val
}

func (p *packet) product() uint64 {
	val := uint64(1)
	for _, sp := range p.subpackets {
		spval := sp.eval()
		val *= spval
	}
	return val
}

func (p *packet) min() uint64 {
	val := p.subpackets[0].eval()
	for _, sp := range p.subpackets[1:] {
		tval := sp.eval()
		if tval < val {
			val = tval
		}
	}
	return val
}

func (p *packet) max() uint64 {
	val := p.subpackets[0].eval()
	for _, sp := range p.subpackets[1:] {
		tval := sp.eval()
		if tval > val {
			val = tval
		}
	}
	return val
}

func (p *packet) gt() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].eval() > p.subpackets[1].eval() {
		return 1
	}
	return 0
}

func (p *packet) lt() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].eval() < p.subpackets[1].eval() {
		return 1
	}
	return 0
}

func (p *packet) eq() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].eval() == p.subpackets[1].eval() {
		return 1
	}
	return 0
}

func (p *packet) eval() uint64 {
	switch int(p.header.typeID) {
	case 4:
		return p.literal()
	case 0: // sum
		return p.sum()
	case 1: // product
		return p.product()
	case 2: // min
		return p.min()
	case 3: // max
		return p.max()
	case 5: // gt
		return p.gt()
	case 6: // lt
		return p.lt()
	case 7: // eq
		return p.eq()
	}

	panic(errors.New(fmt.Sprintf("unknown type id %d", p.header.typeID)))
}

func isPadding(data string) bool {
	for _, c := range data {
		if c != '0' {
			return false
		}
	}
	return true
}

func parse(data string, max int) []*packet {
	packets := make([]*packet, 0, 100)

	for len(data) > 0 && !isPadding(data) {

		p := packet{}
		p.header = newPacketHeader(data[0:6])

		if p.header.typeID == byte(4) {

			i := 6
			for data[i] == '1' {
				i += 5
			}
			p.payload = data[6 : i+5]

		} else {

			lenType := data[6]
			if lenType == '0' {
				p.payload = data[6:22]
			} else {
				p.payload = data[6:18]
			}
			val, _ := strconv.ParseUint(p.payload[1:], 2, 64)
			length := int(val)

			start := 6 + len(p.payload)
			if lenType == '0' {
				end := start + length
				p.subpackets = parse(data[start:end], -1)
			} else {
				p.subpackets = parse(data[start:], length)
			}

		}

		data = data[p.length():]

		packets = append(packets, &p)

		if max > 0 {
			if len(packets) == max {
				break
			}
		}
	}

	return packets
}

func binary(val string) string {
	b := make([]byte, 0, len(val)*4)
	for _, c := range val {
		v, _ := strconv.ParseInt(string(c), 16, 32)
		str := fmt.Sprintf("%04b", v)
		b = append(b, str...)
	}
	return string(b)
}

func getTransmission(filename string) string {
	lines, _ := file.GetLines(filename)
	return binary(lines[0])
}
