package vm

import "bytes"

// ID is a Value representing an address on the blockchain
type ID struct {
	id []byte
}

// NewID creates an ID
func NewID(ab []byte) ID {
	return ID{id: ab}
}

// Compare implements comparison for ID
func (vt ID) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case ID:
		return bytes.Compare(vt.id, other.id), nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt ID) IsScalar() bool {
	return false
}

func (vt ID) String() string {
	return string(vt.id)
}
