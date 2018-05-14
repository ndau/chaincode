package vm

// Constants for time
const (
	EpochStart      = "2018-01-01T00:00:00Z"
	TimestampFormat = "2006-01-02T15:04:05Z"
)

// Value objects are what is managed by the VM
type Value interface {
	Compare(rhs Value) (int, error)
	String() string
}
