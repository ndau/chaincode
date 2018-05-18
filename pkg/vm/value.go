package vm

// ValueType is a special type for the constants that define the different kinds of values
// we can have.
type ValueType int

// Value objects are what is managed by the VM
type Value interface {
	Compare(rhs Value) (int, error)
	IsScalar() bool
	String() string
}
