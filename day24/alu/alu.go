package alu

import (
	"encoding/json"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
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
	W            int64 `json:"w"`
	X            int64 `json:"x"`
	Y            int64 `json:"y"`
	Z            int64 `json:"z"`
	requestInput chan<- bool
	input        <-chan int64
}

func NewALU(requestInput chan<- bool, input <-chan int64) *ALU {
	alu := new(ALU)
	alu.requestInput = requestInput
	alu.input = input
	return alu
}

func (alu *ALU) String() string {
	bytes, _ := json.Marshal(alu)
	return string(bytes)
}

func (alu *ALU) Serialize() string {
	return alu.String()
}

func (alu *ALU) Deserialize(s string) error {
	err := json.Unmarshal([]byte(s), alu)
	if err != nil {
		return err
	}
	return nil
}

func (alu *ALU) Reset() {
	alu.W = 0
	alu.X = 0
	alu.Y = 0
	alu.Z = 0
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
		if v2 == 0 {
			return errors.New("div by 0")
		}
		val := v1 / v2
		alu.set(val1.variable(), val)
	case "mod":
		v1 := val1.value(alu)
		v2 := val2.value(alu)
		if v1 < 0 {
			return errors.New("invalid mod op, val 1 < 0")
		}
		if v2 <= 0 {
			return errors.New("invalid mod op, val 2 <= 0")
		}
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
		return alu.W
	case "x":
		return alu.X
	case "y":
		return alu.Y
	case "z":
		return alu.Z
	}
	panic(errors.New("how did we not match a variable? check the validation"))
}

func (alu *ALU) set(variableName string, value int64) {
	switch variableName {
	case "w":
		alu.W = value
	case "x":
		alu.X = value
	case "y":
		alu.Y = value
	case "z":
		alu.Z = value
	}
}

func RunALU(program []string, inputs []int64) (int64, int64, int64, int64, string, error) {
	return aluRun(program, inputs, `{"w":0,"x":0,"y":0,"z":0}`)
}

func RunALUFromState(program []string, inputs []int64, state string) (int64, int64, int64, int64, string, error) {
	return aluRun(program, inputs, state)
}

func aluRun(runProgram []string, runInputs []int64, initialState string) (int64, int64, int64, int64, string, error) {
	numInput := make(chan int64, 1)
	inputRequest := make(chan bool, 1)

	alu := NewALU(inputRequest, numInput)
	err := alu.Deserialize(initialState)
	if err != nil {
		return 0, 0, 0, 0, "", err
	}

	remainingInputs := runInputs[0:]

	cleanup := func() {
		close(inputRequest)
		close(numInput)
	}

	go func() {
		opened := true
		for opened {
			_, opened = <-inputRequest
			if opened {
				numInput <- remainingInputs[0]
				remainingInputs = remainingInputs[1:]
			}
		}
	}()

	for _, line := range runProgram {
		if len(line) > 0 {
			err := alu.Process(line)
			if err != nil {
				return 0, 0, 0, 0, "", err
			}
		}
	}

	cleanup()

	return alu.W, alu.X, alu.Y, alu.Z, alu.String(), nil
}
