package packet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Reverse(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{`123`, `321`},
		{`1234`, `4321`},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			str := reverse(test.input)
			assert.Equal(t, test.output, str)
		})

	}
}

func Test_CreateLiteralBitString(t *testing.T) {
	tests := []struct {
		input  uint64
		output string
	}{
		{2021, `101111111000101`},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			str := createLiteralBitString(test.input)
			assert.Equal(t, test.output, str)
		})

	}

}

func Test_GetSubpacketSourceStrings(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{`(sum 1 3) (product (sum 3 4) 2)`, []string{`(sum 1 3)`, `(product (sum 3 4) 2)`}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			strs, err := getSubpacketSourceStrings(test.input)
			assert.Nil(t, err)
			assert.Equal(t, len(test.output), len(strs))

			for i := 0; i < len(strs); i++ {
				assert.Equal(t, test.output[i], strs[i])
			}
		})

	}

}

func Test_Compile(t *testing.T) {
	tests := []struct {
		src      string
		compiled string
	}{
		{`(sum 1 3)`, ``},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			compiled, err := Compile(test.src)
			assert.Nil(t, err)
			assert.Equal(t, test.compiled, compiled)
		})

	}

}
