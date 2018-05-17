package vm

import "strings"

var maxStackDepth = 128

// Stack implements the runtime stack for our VM
type Stack struct {
	stack []Value
}

func newStack() *Stack {
	return &Stack{stack: []Value{}}
}

func stackError(s string) error {
	return newRuntimeError("stack " + s)
}

// Clone makes a snapshot copy of a stack
func (st *Stack) Clone() *Stack {
	newstack := make([]Value, len(st.stack))
	for i := range st.stack {
		newstack[i] = st.stack[i]
	}
	return &Stack{stack: newstack}
}

// Depth returns the depth of the stack
func (st *Stack) Depth() int {
	return len(st.stack)
}

// Push puts a value on top of the stack
func (st *Stack) Push(v Value) error {
	if len(st.stack) > maxStackDepth {
		return stackError("overflow")
	}
	st.stack = append(st.stack, v)
	return nil
}

// Get retrieves the item at index n and returns it
func (st *Stack) Get(n int) (Value, error) {
	if len(st.stack) < n {
		return NewNumber(0), stackError("index error")
	}
	last := len(st.stack) - 1
	retval := st.stack[last-n]
	return retval, nil
}

// TopN copies the top N items on the stack and returns a new stack with them
func (st *Stack) TopN(n int) (*Stack, error) {
	if n == 0 {
		return &Stack{[]Value{}}, nil
	}
	if len(st.stack) < n {
		return nil, stackError("index error")
	}
	return &Stack{st.stack[len(st.stack)-n:]}, nil
}

// Peek retrieves the top value and returns it
func (st *Stack) Peek() (Value, error) {
	return st.Get(0)
}

// Pop removes the top value and returns it
func (st *Stack) Pop() (Value, error) {
	if len(st.stack) == 0 {
		return NewNumber(0), stackError("underflow")
	}
	last := len(st.stack) - 1
	retval := st.stack[last]
	st.stack = st.stack[:last]
	return retval, nil
}

// PopAsInt64 retrieves the top entry on the stack as an int64 or errors
func (st *Stack) PopAsInt64() (int64, error) {
	v, err := st.Pop()
	if err != nil {
		return 0, err
	}
	vn, ok := v.(Number)
	if !ok {
		return 0, stackError("top was not number")
	}
	return vn.AsInt64(), nil
}

// PopAsList retrieves the top entry on the stack as a List or errors
func (st *Stack) PopAsList() (List, error) {
	v, err := st.Pop()
	if err != nil {
		return NewList(), err
	}
	l, ok := v.(List)
	if !ok {
		return NewList(), stackError("top was not list")
	}
	return l, nil
}

// PopAsStruct retrieves the top entry on the stack as a Struct or errors
func (st *Stack) PopAsStruct() (Struct, error) {
	v, err := st.Pop()
	if err != nil {
		return NewStruct(), err
	}
	l, ok := v.(Struct)
	if !ok {
		return NewStruct(), stackError("top was not struct")
	}
	return l, nil
}

// PopAt removes the nth value and returns it
func (st *Stack) PopAt(n int) (Value, error) {
	if n == 0 {
		return st.Pop()
	}
	if len(st.stack) <= n {
		return NewNumber(0), stackError("index error")
	}
	ix := len(st.stack) - n - 1
	retval := st.stack[ix]
	st.stack = append(st.stack[:ix], st.stack[ix+1:]...)
	return retval, nil
}

// String renders a stack with one line per value
func (st *Stack) String() string {
	if len(st.stack) == 0 {
		return "|== Empty"
	}
	sa := make([]string, len(st.stack))
	for i := range st.stack {
		sa[i] = "|== " + st.stack[len(st.stack)-i-1].String()
	}
	return strings.Join(sa, "\n")
}
