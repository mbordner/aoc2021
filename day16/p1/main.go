package main

import (
	"aoc2021/common/file"
	"fmt"
	"strconv"
)

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

func main() {
	transmission := getTransmission("../data.txt")

	packets := parse(transmission, -1)

	var verSum func(p *packet) int
	verSum = func(p *packet) int {
		sum := int(p.header.version)
		if p.subpackets != nil {
			for _, sp := range p.subpackets {
				sum += verSum(sp)
			}
		}
		return sum
	}

	fmt.Println(verSum(packets[0]))
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
			val, _ := strconv.ParseInt(p.payload[1:], 2, 32)
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

func getTransmission(filename string) string {
	lines, _ := file.GetLines(filename)

	binary := make([]byte, 0, len(lines[0])*4)
	for _, c := range lines[0] {
		val, _ := strconv.ParseInt(string(c), 16, 32)
		str := fmt.Sprintf("%04b", val)
		binary = append(binary, str...)
	}
	return string(binary)
}
