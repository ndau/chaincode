package vm

import "strings"

const MaxListSize = 1024 // max number of List elements that can result from append or extend

// List maintains a single list object
type List []Value

// NewList creates a new, empty list.
func NewList() List {
	return []Value{}
}

// Compare implements comparison for List
// If lists are different lengths, the shorter
// is "less than" the longer.
// If they are the same length, they are compared
// and the result is the first element that compares nonzero.
func (vt List) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case List:
		if len(vt) < len(other) {
			return -1, nil
		} else if len(vt) > len(other) {
			return 1, nil
		}
		for i := range vt {
			if r, err := vt[i].Compare(other[i]); err != nil || r != 0 {
				return r, err
			}
		}
		return 0, nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

func (vt List) String() string {
	sa := make([]string, len(vt))
	for i, v := range vt {
		sa[i] = v.String()
	}
	return "[" + strings.Join(sa, ", ") + "]"
}

// Len returns the length of a List as an int64
func (vt List) Len() int64 {
	return int64(len(vt))
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
