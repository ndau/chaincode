package vm

import (
	"fmt"
	"strings"
)

// Struct maintains a single struct object; it maintains an array of fields
type Struct struct {
	id     byte
	fields []Value
}

// assert that Struct really is a Value
var _ = Value(Struct{})

// NewStruct creates a new struct with an arbitrary set of fields.
func NewStruct(vs ...Value) Struct {
	return Struct{fields: vs}
}

// Append adds a new field to the end of the Struct and returns it as a new Struct
func (vt Struct) Append(v Value) Struct {
	vt.fields = append(vt.fields, v)
	return vt
}

// Field retrieves the field at a given index
func (vt Struct) Field(ix int) (Value, error) {
	if ix >= len(vt.fields) || ix < 0 {
		return NewNumber(0), ValueError{"invalid field index"}
	}
	return vt.fields[ix], nil
}

// Less implements comparison for Struct
// If structs have different IDs, or rhs is not a Struct, errors.
// If they are the same ID, they are compared field by field
// and the result is the first element that compares nonzero.
// If the iteration runs off the end, the shorter struct is less.
func (vt Struct) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case Struct:
		for i := 0; true; i++ {
			// if the structs have compared equal so far (which they have since we got here)
			// and v2 runs off the end, then the result is definitely false
			v2, err := other.Field(i)
			if err != nil {
				return false, nil
			}
			// if v1 runs off the end first, then the result is true
			v1, err := vt.Field(i)
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
func (vt Struct) IsScalar() bool {
	return false
}

func (vt Struct) String() string {
	sa := make([]string, len(vt.fields))
	for i, v := range vt.fields {
		sa[i] = v.String()
	}
	return fmt.Sprintf("str(%d)[%s]", vt.id, strings.Join(sa, ", "))
}

// IsTrue indicates if this Value evaluates to true
func (vt Struct) IsTrue() bool {
	return false
}
