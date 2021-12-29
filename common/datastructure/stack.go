package datastructure

type Stack struct {
	vals []interface{}
}

func NewStack(cap int) *Stack {
	s := new(Stack)
	s.vals = make([]interface{}, 0, cap)
	return s
}

func (s *Stack) Push(val interface{}) {
	s.vals = append(s.vals, val)
}

func (s *Stack) Pop() interface{} {
	if len(s.vals) > 0 {
		val := s.vals[len(s.vals)-1]
		s.vals = s.vals[0 : len(s.vals)-1]
		return val
	}
	return nil
}

func (s *Stack) Peek() interface{} {
	if len(s.vals) > 0 {
		return s.vals[len(s.vals)-1]
	}
	return nil
}
