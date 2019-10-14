package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

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
	return fmt.Sprintf("%s [pc=%d]", e.msg, e.pc)
}
