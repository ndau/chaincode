package vm

// Value objects are what is managed by the VM
type Value interface {
	Compare(rhs Value) (int, error)
	IsScalar() bool
	String() string
}
