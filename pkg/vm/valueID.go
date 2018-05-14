package vm

import "strconv"

// ID is a Value representing an address on the blockchain
type ID struct {
	id uint64
}

// NewID creates an ID
func NewID(n uint64) ID {
	return ID{n}
}

// Compare implements comparison for ID
func (vt ID) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case ID:
		if vt.id < other.id {
			return -1, nil
		} else if vt.id > other.id {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

func (vt ID) String() string {
	return strconv.FormatUint(vt.id, 16)
}
