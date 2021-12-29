package packet

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Decompile(t *testing.T) {
	tests := []struct {
		input  string
		output uint64
		source string
	}{
		{"324A", 42, `42`},
		{"302", 1, `1`},
		{"3C0080C928C928", 1, `(eq 42 42)`},
		{"200058C0983", 4, `(sum 1 3)`},
		{"C200B40A82", 3, `(sum 1 2)`},
		{"04005AC33890", 54, `(product 6 9)`},
		{"880086C3E88112", 7, `(min 7 8 9)`},
		{"CE00C43D881120", 9, `(max 7 8 9)`},
		{"D8005AC2A8F0", 1, `(lt 5 15)`},
		{"F600BC2D8F", 0, `(gt 5 15)`},
		{"9C005AC2F8F0", 0, `(eq 5 15)`},
		{"9C0141080250320F1802104A08", 1, `(eq (sum 1 3) (product 2 2))`},
		{"3C01608001630260C90016304608", 1, `(eq (sum 1 3) (product 2 2))`},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("Test Case %d", i+1), func(t *testing.T) {
			p := Decompile(test.input)
			assert.Equal(t, test.output, p.Eval())
			assert.Equal(t, test.source, p.String())
		})

	}

}
