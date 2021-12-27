package alu

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestALU_Process(t *testing.T) {

	cases := []struct {
		program  []string
		inputs   []int64
		expected string
	}{
		{
			program:  breakInput(input1()),
			inputs:   []int64{3},
			expected: `{w: 0, x: -3, y: 0, z: 0}`,
		},
		{
			program:  breakInput(input1()),
			inputs:   []int64{-3},
			expected: `{w: 0, x: 3, y: 0, z: 0}`,
		},
		{
			program:  breakInput(input2()),
			inputs:   []int64{3, 9},
			expected: `{w: 0, x: 9, y: 0, z: 1}`,
		},
		{
			program:  breakInput(input2()),
			inputs:   []int64{2, 9},
			expected: `{w: 0, x: 9, y: 0, z: 0}`,
		},
		{
			program:  []string{`inp z`, `inp y`, `div z y`},
			inputs:   []int64{9, 3},
			expected: `{w: 0, x: 0, y: 3, z: 3}`,
		},
		{
			program:  []string{`inp z`, `inp y`, `div z y`},
			inputs:   []int64{9, 2},
			expected: `{w: 0, x: 0, y: 2, z: 4}`,
		},
		{
			program:  []string{`inp z`, `inp y`, `mod z y`},
			inputs:   []int64{9, 2},
			expected: `{w: 0, x: 0, y: 2, z: 1}`,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {

			numInput := make(chan int64, 1)
			inputRequest := make(chan bool, 1)

			alu := NewALU(inputRequest, numInput)

			inputs := tc.inputs[0:]

			go func() {
				opened := true
				for opened {
					_, opened = <-inputRequest
					if opened {
						numInput <- inputs[0]
						inputs = inputs[1:]
					}
				}
			}()

			for _, line := range tc.program {
				if len(line) > 0 {
					err := alu.Process(line)
					assert.Nil(t, err)
				}
			}

			close(inputRequest)
			close(numInput)

			assert.Equal(t, tc.expected, alu.String())

		})
	}
}

func breakInput(s string) []string {
	return strings.Split(s, "\n")
}

func input1() string {
	return `inp x
mul x -1`
}

func input2() string {
	return `inp z
inp x
mul z 3
eql z x
`
}

func input3() string {
	return `inp w
add z w
mod z 2
div w 2
add y w
mod y 2
div w 2
add x w
mod x 2
div w 2
mod w 2`
}
