package vm

import (
	"strconv"
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

// ValueError is the error type when value conflicts arise
type ValueError struct {
	msg string
}

func (e ValueError) Error() string {
	return e.msg
}

// Value objects are what is managed by the VM
type Value interface {
	String() string
}

// A Number is a type of Value representing an int64
type Number struct {
	v int64
}

func newNumber(n int64) Number {
	return Number{n}
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

func newTimestamp(n uint64) Timestamp {
	return Timestamp{n}
}

func (vt Timestamp) String() string {
	return strconv.FormatUint(vt.t, 16)
}

// ID is a Value representing an address on the blockchain
type ID struct {
	id uint64
}

func newID(n uint64) ID {
	return ID{n}
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
