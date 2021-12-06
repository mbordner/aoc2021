package expression

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

// https://www.engr.mun.ca/~theo/Misc/exp_parsing.htm

var (
	reSpace    = regexp.MustCompile(`\s`)
	reOperator = regexp.MustCompile(`\+|\*`)
	reDigits   = regexp.MustCompile(`\d+`)
)

type Parser struct {
	operators    []*Operator
	operands     []interface{}
	opPrecedence Precedence
	start        int
	end          int
	expr         string
}

func (p *Parser) E() {
	p.P()
	n, err := p.next()
	if err != nil {
		panic(err)
	}
	for reOperator.MatchString(n) && IsBinary(n) {
		p.pushOperator(n)
		p.consume()

		p.P()

		n, err = p.next()
		if err != nil {
			panic(err)
		}
	}
	for p.operators[len(p.operators)-1] != nil {
		p.popOperator()
	}
}

func (p *Parser) P() {
	n, err := p.next()
	if err != nil {
		panic(err)
	}
	if reDigits.MatchString(n) {
		v, _ := strconv.ParseInt(n, 10, 64)
		p.operands = append(p.operands, v)
		p.consume()
	} else if n == "(" {
		p.consume()
		p.operators = append(p.operators, nil)
		p.E()
		n, err = p.next()
		if err != nil {
			panic(err)
		}
		if n != ")" {
			panic(errors.New("expected )"))
		}
		p.consume()
		p.operators = p.operators[0 : len(p.operators)-1]
	} else {
		panic(errors.New("error parsing expression"))
	}
}

func (p *Parser) popOperator() {
	op := p.operators[len(p.operators)-1]
	if IsBinary(op.op) {
		p.operators = p.operators[0 : len(p.operators)-1]
		op.right = p.operands[len(p.operands)-1]
		op.left = p.operands[len(p.operands)-2]
		p.operands = p.operands[0 : len(p.operands)-2]
		p.operands = append(p.operands, op)
	}
}

func (p *Parser) pushOperator(op string) {
	for p.operators[len(p.operators)-1] != nil && p.opPrecedence(op, p.operators[len(p.operators)-1].op) <= 0 {
		p.popOperator()
	}
	o := Operator{}
	o.op = op
	p.operators = append(p.operators, &o)
}

func (p *Parser) next() (string, error) {

	for p.start < len(p.expr) && reSpace.MatchString(string(p.expr[p.start])) {
		p.start++
	}

	if p.start == len(p.expr) {
		return "", nil
	}

	p.end = p.start

	if reDigits.MatchString(string(p.expr[p.start])) {
		for p.end < len(p.expr) && reDigits.MatchString(string(p.expr[p.end])) {
			p.end++
		}
	} else if p.expr[p.start] == '(' || p.expr[p.start] == ')' {
		p.end++
	} else if reOperator.MatchString(string(p.expr[p.start])) {
		p.end++
	} else {
		return "", errors.New("unexpected token")
	}

	return p.expr[p.start:p.end], nil
}

func (p *Parser) consume() {
	p.start = p.end
}

func (p *Parser) Eval() int64 {
	return p.operands[0].(*Operator).Eval()
}

func NewParser(expr string, precedence Precedence) *Parser {
	p := Parser{}
	p.opPrecedence = precedence
	p.expr = expr
	p.operators = make([]*Operator, 0, 20)
	p.operands = make([]interface{}, 0, 20)

	p.operators = append(p.operators, nil)
	p.E()

	return &p
}
