package vm

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import "strconv"

// A Number is a type of Value representing an int64
type Number struct {
	v int64
}

// assert that Number really is a Value
var _ = Value(Number{})

// NewNumber creates a Number object out of an int64
func NewNumber(n int64) Number {
	return Number{n}
}

// NewTrue returns a number that is considered to be true
func NewTrue() Number {
	return Number{-1}
}

// NewFalse returns a number that is considered to be false
func NewFalse() Number {
	return Number{0}
}

// Equal implements equality testing for Number
func (vt Number) Equal(rhs Value) bool {
	switch other := rhs.(type) {
	case Number:
		return vt.v == other.v
	default:
		return false
	}
}

// Less implements comparison for Number
func (vt Number) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case Number:
		return vt.v < other.v, nil
	default:
		return false, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt Number) IsScalar() bool {
	return true
}

func (vt Number) String() string {
	return strconv.FormatInt(vt.v, 10)
}

// IsTrue indicates if this Value evaluates to true
func (vt Number) IsTrue() bool {
	return vt.v != 0
}

// AsInt64 allows retrieving the contents of a Number object as an int64
//
// implements Numeric
func (vt Number) AsInt64() int64 {
	return vt.v
}

var _ Numeric = (*Number)(nil)
