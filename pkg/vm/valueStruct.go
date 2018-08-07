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

// Compare implements comparison for Struct
// If structs have different IDs, or rhs is not a Struct, errors.
// If they are the same ID, they are compared field by field
// and the result is the first element that compares nonzero.
func (vt Struct) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case Struct:
		if vt.id != other.id {
			return 0, ValueError{"comparing different structs"}
		}
		for i := range vt.fields {
			if r, err := vt.fields[i].Compare(other.fields[i]); err != nil || r != 0 {
				return r, err
			}
		}
		return 0, nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
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
