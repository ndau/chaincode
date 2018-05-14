package vm

import "fmt"

// ValidationError is returned when the code is invalid and cannot be loaded or run
type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}

// ValueError is the error type when value conflicts arise
type ValueError struct {
	msg string
}

func (e ValueError) Error() string {
	return e.msg
}

// RuntimeError is returned when the vm encounters an error during execution
type RuntimeError struct {
	pc  int
	msg string
}

// PC sets the program counter value for an error
func (e RuntimeError) PC(pc int) RuntimeError {
	e.pc = pc
	return e
}

func newRuntimeError(s string) error {
	return RuntimeError{pc: -1, msg: s}
}

func wrapRuntimeError(e error) RuntimeError {
	if rte, ok := e.(RuntimeError); ok {
		return rte
	}
	return RuntimeError{pc: -1, msg: e.Error()}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[pc=%d] %s", e.pc, e.msg)
}
