package expression

import "github.com/pkg/errors"

type Precedence func(op1, op2 string) int

type Operator struct {
	op    string
	left  interface{}
	right interface{}
}

func (o *Operator) Eval() int64 {
	var l, r int64
	switch o.left.(type) {
	case int64:
		l = o.left.(int64)
	case *Operator:
		l = o.left.(*Operator).Eval()
	}
	switch o.right.(type) {
	case int64:
		r = o.right.(int64)
	case *Operator:
		r = o.right.(*Operator).Eval()
	}
	switch o.op {
	case "+":
		return l + r
	case "*":
		return l * r
	}
	panic(errors.New("unknown operator"))
}

func CompareOperator(op1, op2 string) int {
	return 0
}

func IsBinary(op string) bool {
	return true
}
