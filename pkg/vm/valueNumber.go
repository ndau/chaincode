package vm

import "strconv"

// A Number is a type of Value representing an int64
type Number struct {
	v int64
}

// NewNumber creates a Number object out of an int64
func NewNumber(n int64) Number {
	return Number{n}
}

// Compare implements comparison for Number
func (vt Number) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case Number:
		if vt.v < other.v {
			return -1, nil
		} else if vt.v > other.v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

func (vt Number) String() string {
	return strconv.FormatInt(vt.v, 10)
}

// AsInt64 allows retrieving the contents of a Number object as an int64
func (vt Number) AsInt64() int64 {
	return vt.v
}
