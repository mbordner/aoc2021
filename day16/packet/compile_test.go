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
		{1, `00001`},
		{0, `00000`},
		{16, `1000100000`},
		{17, `1000100001`},
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
		{`42`, `324A`},
		{`1`, `302`},
		{`(eq 42 42)`, `3C0080C928C928`},
		{`(sum 1 3)`, `200058C0983`},
		{`(eq (sum 1 3) (product 2 2))`, `3C01608001630260C90016304608`},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			p, hex, err := Compile(test.src)
			assert.NotNil(t, p)
			assert.Nil(t, err)
			assert.Equal(t, test.compiled, hex)
		})

	}

}
