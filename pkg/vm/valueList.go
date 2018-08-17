package vm

import (
	"strings"
)

// MaxListSize is the max number of List elements that can result from append or extend
const MaxListSize = 1024

// List maintains a single list object
type List []Value

// assert that List really is a Value
var _ = Value(List{})

// NewList creates a new, empty list.
func NewList(vs ...Value) List {
	if vs == nil {
		return make(List, 0)
	}
	return vs
}

// Less implements comparison for List
// The lists are compared on a per-element basis
// and the result is the result of the first element
// that is not equal to its counterpart.
// If lists are different lengths but equal for comparative
// lengths, the shorter is less than the longer.
func (vt List) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case List:
		for i := int64(0); true; i++ {
			// if the lists have compared equal so far (which they have since we got here)
			// and v2 runs off the end, then the result is definitely false
			v2, err := other.Index(i)
			if err != nil {
				return false, nil
			}
			// if v1 runs off the end first, then the result is true
			v1, err := vt.Index(i)
			if err != nil {
				return true, nil
			}
			// if v1 < v2 errors return the error
			r1, err := v1.Less(v2)
			if err != nil {
				return false, err
			}
			// if v1 < v2 return true
			if r1 {
				return true, nil
			}
			// if v1 > v2 return false, otherwise go around again
			if r2, _ := v2.Less(v1); r2 {
				return false, nil
			}
		}
	default:
		return false, ValueError{"comparing incompatible types"}
	}
	// this is here because go's escape analysis is failing
	panic("List: can't happen")
}

// IsScalar indicates if this Value is a scalar value type
func (vt List) IsScalar() bool {
	return false
}

func (vt List) String() string {
	sa := make([]string, len(vt))
	for i, v := range vt {
		sa[i] = v.String()
	}
	return "[" + strings.Join(sa, ", ") + "]"
}

// IsTrue indicates if this Value evaluates to true
func (vt List) IsTrue() bool {
	return false
}

// Len returns the length of a List as an int64
func (vt List) Len() int64 {
	return int64(len(vt))
}

// Index returns the value at the given index, or error
func (vt List) Index(n int64) (Value, error) {
	if n >= vt.Len() || n < -vt.Len() {
		return nil, ValueError{"list index out of bounds"}
	}
	if n < 0 {
		return vt[int(vt.Len()+n)], nil
	}
	return vt[n], nil
}

// Append adds a new Value to the end of a list
func (vt List) Append(v Value) List {
	return append(vt, v)
}

// Extend generates a new List by appending one List to the end of another
func (vt List) Extend(other List) List {
	return append(vt, other...)
}

// Map applies a function to each element of the list and returns a List of the results
func (vt List) Map(f func(Value) (Value, error)) (List, error) {
	result := NewList()
	for _, v := range vt {
		r, err := f(v)
		if err != nil {
			return result, err
		}
		result = result.Append(r)
	}
	return result, nil
}

// Reduce applies a function to each element of the list and returns an aggregated result
func (vt List) Reduce(f func(prev, item Value) Value, init Value) Value {
	result := init
	for _, v := range vt {
		result = f(result, v)
	}
	return result
}
