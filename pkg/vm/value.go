package vm

// Value objects are what is managed by the VM
type Value interface {
	Equal(rhs Value) bool
	Less(rhs Value) (bool, error)
	IsScalar() bool
	String() string
	IsTrue() bool
}
