package main

import (
	"aoc2021/common/file"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	reDigit    = regexp.MustCompile(`-?\d+`)
	reVariable = regexp.MustCompile(`[w-z]`)
)

type aluValue struct {
	variableName string
	val          int64
}

func newALUValue(expr string) (*aluValue, error) {
	v := new(aluValue)
	if reVariable.MatchString(expr) {
		v.variableName = expr
	} else if reDigit.MatchString(expr) {
		val, err := strconv.ParseInt(expr, 10, 64)
		if err != nil {
			return nil, err
		}
		v.val = val
	} else {
		return nil, errors.New("invalid expression")
	}
	return v, nil
}

func (v *aluValue) value(alu *ALU) int64 {
	if v.variableName != "" {
		return alu.get(v.variableName)
	}
	return v.val
}

func (v *aluValue) variable() string {
	return v.variableName
}

type ALU struct {
	w            int64
	x            int64
	y            int64
	z            int64
	requestInput chan<- bool
	input        <-chan int64
}

func NewALU(requestInput chan<- bool, input <-chan int64) *ALU {
	alu := new(ALU)
	alu.requestInput = requestInput
	alu.input = input
	return alu
}

func (alu *ALU) W() int64 {
	return alu.w
}

func (alu *ALU) X() int64 {
	return alu.x
}

func (alu *ALU) Y() int64 {
	return alu.y
}

func (alu *ALU) Z() int64 {
	return alu.z
}

func (alu *ALU) String() string {
	return fmt.Sprintf("{w: %d, x: %d, y: %d, z: %d}", alu.w, alu.x, alu.y, alu.z)
}

func (alu *ALU) Process(instruction string) error {
	instruction = strings.TrimSpace(instruction)
	if len(instruction) == 0 {
		return nil
	}
	tokens := strings.Split(instruction, " ")
	instr := tokens[0]
	var val1, val2 *aluValue
	var err error

	val1, err = newALUValue(tokens[1])
	if err != nil {
		return err
	}
	if len(tokens) > 2 {
		val2, err = newALUValue(tokens[2])
		if err != nil {
			return err
		}
	}

	if val1.variable() == "" {
		return errors.New("instruction must have variableName as first param")
	}

	switch instr {
	case "inp":
		alu.requestInput <- true
		val := <-alu.input
		alu.set(val1.variable(), val)
	case "add":
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		val := v1 + v2
		alu.set(val1.variable(), val)
	case "mul":
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		val := v1 * v2
		alu.set(val1.variable(), val)
	case "div":
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		val := v1 / v2
		alu.set(val1.variable(), val)
	case "mod":
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		val := v1 % v2
		alu.set(val1.variable(), val)
	case "eql":
		val := int64(0)
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		if v1 == v2 {
			val = int64(1)
		}
		alu.set(val1.variable(), val)
	default:
		return errors.New("unknown instruction type")
	}

	return nil
}

func (alu *ALU) get(variableName string) int64 {
	switch variableName {
	case "w":
		return alu.w
	case "x":
		return alu.x
	case "y":
		return alu.y
	case "z":
		return alu.z
	}
	panic(errors.New("how did we not match a variable? check the validation"))
}

func (alu *ALU) set(variableName string, value int64) {
	switch variableName {
	case "w":
		alu.w = value
	case "x":
		alu.x = value
	case "y":
		alu.y = value
	case "z":
		alu.z = value
	}
}

func main() {

	numInput := make(chan int64, 1)
	inputRequest := make(chan bool, 1)

	alu := NewALU(inputRequest, numInput)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	cleanup := func(code int) {
		close(signals)
		close(inputRequest)
		close(numInput)
		os.Exit(code)
	}

	go func() {
		for {
			select {
			case <-signals:
				cleanup(1)
			case <-inputRequest:
				delim := '\n'
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("input> ")
				text, err := reader.ReadString(byte(delim))
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "%v", err)
					cleanup(-3)
				}

				tokens := strings.Split(text, string(delim))
				if len(tokens) > 0 {
					text = tokens[0]
				}

				if text == "exit" || text == "quit" {
					fmt.Println("exiting, goodbye!")
					cleanup(0)
				} else {
					val, err := strconv.ParseInt(text, 10, 64)
					if err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "error parsing value: %s", err)
					}
					numInput <- val
				}
			}

		}
	}()

	lines, err := file.GetLines(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Can not read file %s", os.Args[1])
		cleanup(-1)
	}

	for num, line := range lines {
		if len(line) > 0 {
			fmt.Println("processing instruction (", num, "): ", line)
			err = alu.Process(line)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v", err)
				cleanup(-2)
			}
		}
	}

	fmt.Println(alu.String())
	cleanup(0)

}
