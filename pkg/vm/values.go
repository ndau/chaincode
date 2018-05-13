package vm

import (
	"strconv"
	"time"
)

// ValueType is a byte indicator for a given
// type ValueType byte

// const (
// 	VtNumber    ValueType = iota
// 	VtID        ValueType = iota
// 	VtTimestamp ValueType = iota
// 	VtList      ValueType = iota
// 	VtStruct    ValueType = iota
// )

// Constants for time
const (
	EpochStart      = "2018-01-01T00:00:00Z"
	TimestampFormat = "2006-01-02T15:04:05Z"
)

// Some masking values
const (
	ByteMask byte = 0xFF
	HighBit  byte = 0x80
)

// ToBytesU returns an array of 8 bytes encoding n as a uint in little-endian form
func ToBytesU(n uint64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// ToBytes returns an array of 8 bytes encoding n as a signed value in little-endian form
func ToBytes(n int64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// ValueError is the error type when value conflicts arise
type ValueError struct {
	msg string
}

func (e ValueError) Error() string {
	return e.msg
}

// Value objects are what is managed by the VM
type Value interface {
	Compare(rhs Value) (int, error)
	String() string
}

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

// Timestamp is a Value type representing duration since the epoch
type Timestamp struct {
	t uint64
}

// NewTimestamp creates a timestamp from a uint64 representation of one
func NewTimestamp(n uint64) Timestamp {
	return Timestamp{n}
}

// ParseTimestamp creates a timestamp from an ISO-3933 string
func ParseTimestamp(s string) (Timestamp, error) {
	epoch, err := time.Parse(TimestampFormat, EpochStart)
	if err != nil {
		panic("Epoch isn't a valid timestamp!")
	}
	ts, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return NewTimestamp(0), err
	}
	// durations are in nanoseconds but we want microseconds
	uSec := uint64(ts.Sub(epoch).Nanoseconds() / 1000)
	return NewTimestamp(uSec), nil
}

// Compare implements comparison for Timestamp
func (vt Timestamp) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case Timestamp:
		if vt.t < other.t {
			return -1, nil
		} else if vt.t > other.t {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

func (vt Timestamp) String() string {
	return strconv.FormatUint(vt.t, 16)
}

// T returns the timestamp as a uint64 duration in uSec since the start of epoch.
func (vt Timestamp) T() uint64 {
	return vt.t
}

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

// // List maintains a single list object
// type ListList struct {
// 	list []Value
// }

// type List struct {
// 	index int
// }

// func newList() List {
// 	return []Value{}
// }

// func (vt List) String() string {
// 	sa := make([]string, len(vt))
// 	for i, v := range vt {
// 		sa[i] = v.String()
// 	}
// 	return "[" + strings.Join(sa, ", ") + "]"
// }
