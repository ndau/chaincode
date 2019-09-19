package vm

// Numeric types can be expressed as integers
type Numeric interface {
	AsInt64() int64
}
