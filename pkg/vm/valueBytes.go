package vm

import "bytes"

// Bytes is a Value representing an address on the blockchain
type Bytes struct {
	b []byte
}

// assert that Bytes really is a Value
var _ = Value(&Bytes{})

// NewBytes creates an Bytes
func NewBytes(ab []byte) *Bytes {
	return &Bytes{b: ab}
}

// Less implements comparison for Bytes
func (vt *Bytes) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case *Bytes:
		return bytes.Compare(vt.b, other.b) < 0, nil
	default:
		return false, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt *Bytes) IsScalar() bool {
	return true
}

func (vt *Bytes) String() string {
	return string(vt.b)
}

// IsTrue indicates if this Value evaluates to true
func (vt *Bytes) IsTrue() bool {
	return false
}
