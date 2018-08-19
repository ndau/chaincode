package vm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/bitset256"
)

// Struct maintains a single struct object; it maintains a map of byte ids to fields
type Struct struct {
	validFields *bitset256.Bitset256
	fields      map[byte]Value
}

// assert that Struct really is a Value
var _ = Value(&Struct{})

// NewStruct creates a new, empty struct.
func NewStruct() *Struct {
	return &Struct{
		validFields: bitset256.New(),
		fields:      make(map[byte]Value),
	}
}

// NewTestStruct creates a new struct with an arbitrary list of fields.
// The fields will be created with an index in order beginning at 0.
// This is really only intended for testing.
func NewTestStruct(vs ...Value) *Struct {
	st := NewStruct()
	for i, v := range vs {
		st.Set(byte(i), v)
	}
	return st
}

// Set assigns a value to a field at index ix and returns it.
func (vt *Struct) Set(ix byte, v Value) *Struct {
	vt.validFields.Set(ix)
	vt.fields[ix] = v
	return vt
}

// SafeSet assigns a value to a field at index ix, only if the
// there is no field already at that index.
func (vt *Struct) SafeSet(ix byte, v Value) (*Struct, error) {
	if vt.validFields.Get(ix) {
		return nil, errors.New("attempt to overwrite existing struct member")
	}
	return vt.Set(ix, v), nil
}

// Get retrieves the field at a given index
func (vt *Struct) Get(ix byte) (Value, error) {
	f, ok := vt.fields[ix]
	if !ok {
		return NewNumber(0), ValueError{"invalid field index"}
	}
	return f, nil
}

// IsCompatible returns true if the other struct list of validFields
// is equal to the receiver's list.
func (vt *Struct) IsCompatible(other *Struct) bool {
	return vt.validFields.Equals(other.validFields)
}

// Equal implements comparison for Struct. If rhs is not a Struct, errors. If
// the two structs have different values for validFields, then the result is
// false. If they have the same field set, they are compared field by field in
// numeric order and the result is the result from the first element that is not
// equal to its counterpart.
func (vt *Struct) Equal(rhs Value) bool {
	switch other := rhs.(type) {
	case *Struct:
		if !vt.IsCompatible(other) {
			return false
		}
		fieldIDs := vt.validFields.Indices()
		for _, ix := range fieldIDs {
			// we know that the structs both have the same field IDs so we're
			// safe in ignoring errors
			f1 := vt.fields[ix]
			f2 := other.fields[ix]
			if !f1.Equal(f2) {
				return false
			}
		}
		// if we get here, the two structs were equal, so therefore not less
		return true
	default:
		return false
	}
}

// Less implements comparison for Struct. If rhs is not a Struct, errors. If the
// two structs have different values for validFields, then the result is the
// result of comparing the new validFields objects. If they have the same field
// set, they are compared field by field in numeric order and the result is the
// result from the first element that is not equal to its counterpart.
func (vt *Struct) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case *Struct:
		if !vt.IsCompatible(other) {
			return vt.validFields.Less(other.validFields), nil
		}
		fieldIDs := vt.validFields.Indices()
		for _, ix := range fieldIDs {
			// we know that the structs both have the same field IDs so we're
			// safe in ignoring errors (any type errors at the field level will
			// be caught by Less).
			f1 := vt.fields[ix]
			f2 := other.fields[ix]
			if !f1.Equal(f2) {
				return f1.Less(f2)
			}
		}
		// if we get here, the two structs were equal, so therefore not less
		return false, nil
	default:
		return false, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt *Struct) IsScalar() bool {
	return false
}

func (vt *Struct) String() string {
	sa := make([]string, len(vt.fields))
	i := 0
	for _, k := range vt.validFields.Indices() {
		sa[i] = fmt.Sprintf("%d: %s", k, vt.fields[k].String())
		i++
	}
	return fmt.Sprintf("struct{%s}", strings.Join(sa, ", "))
}

// IsTrue indicates if this Value evaluates to true
func (vt *Struct) IsTrue() bool {
	return false
}

// Len returns the number of fields in this struct
func (vt *Struct) Len() int {
	return len(vt.fields)
}

// Indices returns the list of indices that are defined in this struct.
func (vt *Struct) Indices() []byte {
	return vt.validFields.Indices()
}
