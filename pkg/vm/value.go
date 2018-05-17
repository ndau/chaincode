package vm

// ValueType is a special type for the constants that define the different kinds of values
// we can have.
type ValueType int

// These are ValueType constants
const (
	NumberT    ValueType = iota
	IDT        ValueType = iota
	TimestampT ValueType = iota
	ListT      ValueType = iota
	StructT    ValueType = iota
)

// Value objects are what is managed by the VM
type Value interface {
	Compare(rhs Value) (int, error)
	IsScalar() bool
	String() string
}
