package packet

import (
	"errors"
	"fmt"
	"strconv"
)

type Operation int

const (
	Sum     Operation = iota // 0
	Product                  // 1
	Min                      // 2
	Max                      // 3
	Literal                  // 4
	Gt                       // 5
	Lt                       // 6
	Eq                       // 7
)

const (
	SumStr     = "sum"
	ProductStr = "product"
	MinStr     = "min"
	MaxStr     = "max"
	GtStr      = "gt"
	LtStr      = "lt"
	EqStr      = "eq"
)

var (
	opStrMap = map[Operation]string{
		Sum:     SumStr,
		Product: ProductStr,
		Min:     MinStr,
		Max:     MaxStr,
		Literal: "",
		Gt:      GtStr,
		Lt:      LtStr,
		Eq:      EqStr,
	}
)

func (o Operation) String() string {
	return opStrMap[o]
}

func OperationFromString(s string) (Operation, error) {
	for k, v := range opStrMap {
		if v == s {
			return k, nil
		}
	}
	return Operation(0), errors.New("unknown operator")
}

type packetHeader struct {
	version byte
	typeID  byte
	payload string
}

func newPacketHeaderFromPayload(payload string) *packetHeader {
	ph := new(packetHeader)
	ph.payload = payload

	val, _ := strconv.ParseInt(payload[0:3], 2, 8)
	ph.version = byte(val)

	val, _ = strconv.ParseInt(payload[3:], 2, 8)
	ph.typeID = byte(val)

	return ph
}

func newPacketHeaderFromValues(id Operation, ver int) *packetHeader {
	ph := new(packetHeader)

	ph.version = byte(ver)
	ph.payload = fmt.Sprintf("%03b", ph.version)

	ph.typeID = byte(id)
	ph.payload += fmt.Sprintf("%03b", ph.typeID)

	return ph
}

type Packet struct {
	header     *packetHeader
	payload    string
	subpackets []*Packet
}

func Decompile(data string) *Packet {
	ps := parse(binary(data), 1)
	return ps[0]
}

func (p *Packet) prettyString(indent string) string {
	str := indent + "("

	if Operation(p.header.typeID) == Literal {
		str += fmt.Sprintf("%d", p.literal()) + ")\n"
	} else {
		str += Operation(p.header.typeID).String() + "\n"

		if p.subpackets != nil {
			for _, sp := range p.subpackets {
				str += sp.prettyString(indent + "   ")
			}
		}

		str += indent + ")\n"
	}

	return str
}

func (p *Packet) String() string {
	if Operation(p.header.typeID) == Literal {
		return fmt.Sprintf("%d", p.literal())
	} else {
		str := "(" + Operation(p.header.typeID).String()

		if p.subpackets != nil {
			for _, sp := range p.subpackets {
				str += " " + sp.String()
			}
		}

		str += ")"
		return str
	}
}

func (p *Packet) PrettyString() string {
	return p.prettyString("")
}

func (p *Packet) length() int {
	l := len(p.header.payload) + len(p.payload)
	if p.subpackets != nil {
		for _, sp := range p.subpackets {
			l += sp.length()
		}
	}
	return l
}

func (p *Packet) literal() uint64 {
	value := make([]byte, 0, len(p.payload))
	for i := 0; i < len(p.payload); i += 5 {
		value = append(value, p.payload[i+1:i+5]...)
	}
	val, _ := strconv.ParseUint(string(value), 2, 64)
	return val
}

func (p *Packet) sum() uint64 {
	if len(p.subpackets) < 1 {
		panic(errors.New("supposed to have at least one subpacket"))
	}
	val := uint64(0)
	for _, sp := range p.subpackets {
		spval := sp.Eval()
		val += spval
	}
	return val
}

func (p *Packet) product() uint64 {
	if len(p.subpackets) < 1 {
		panic(errors.New("supposed to have at least one subpacket"))
	}
	val := uint64(1)
	for _, sp := range p.subpackets {
		spval := sp.Eval()
		val *= spval
	}
	return val
}

func (p *Packet) min() uint64 {
	if len(p.subpackets) < 1 {
		panic(errors.New("supposed to have at least one subpacket"))
	}
	val := p.subpackets[0].Eval()
	for _, sp := range p.subpackets[1:] {
		tval := sp.Eval()
		if tval < val {
			val = tval
		}
	}
	return val
}

func (p *Packet) max() uint64 {
	if len(p.subpackets) < 1 {
		panic(errors.New("supposed to have at least one subpacket"))
	}
	val := p.subpackets[0].Eval()
	for _, sp := range p.subpackets[1:] {
		tval := sp.Eval()
		if tval > val {
			val = tval
		}
	}
	return val
}

func (p *Packet) gt() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].Eval() > p.subpackets[1].Eval() {
		return 1
	}
	return 0
}

func (p *Packet) lt() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].Eval() < p.subpackets[1].Eval() {
		return 1
	}
	return 0
}

func (p *Packet) eq() uint64 {
	if len(p.subpackets) != 2 {
		panic(errors.New("supposed to only have 2 subpackets"))
	}
	if p.subpackets[0].Eval() == p.subpackets[1].Eval() {
		return 1
	}
	return 0
}

func (p *Packet) Eval() uint64 {
	switch Operation(p.header.typeID) {
	case Literal:
		return p.literal()
	case Sum:
		return p.sum()
	case Product:
		return p.product()
	case Min:
		return p.min()
	case Max:
		return p.max()
	case Gt:
		return p.gt()
	case Lt:
		return p.lt()
	case Eq:
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

func parse(data string, max int) []*Packet {
	packets := make([]*Packet, 0, 100)

	for len(data) > 0 && !isPadding(data) {

		p := Packet{}
		p.header = newPacketHeaderFromPayload(data[0:6])

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
