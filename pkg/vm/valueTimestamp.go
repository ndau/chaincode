package vm

import (
	"time"

	"github.com/oneiro-ndev/ndaumath/pkg/types"
)

// Timestamp is a Value type representing duration since the epoch
type Timestamp struct {
	t types.Timestamp
}

// assert that Timestamp really is a Value
var _ = Value(Timestamp{})

// NewTimestampFromInt creates a timestamp from an int64 representation of one
func NewTimestampFromInt(n int64) Timestamp {
	return Timestamp{types.Timestamp(n)}
}

// NewTimestamp returns a timestamp taken from a ndaumath/pkg/types.Timestamp struct
func NewTimestamp(t types.Timestamp) Timestamp {
	return Timestamp{t: t}
}

// NewTimestampFromTime returns a timestamp taken from a time.Time struct in Go.
func NewTimestampFromTime(t time.Time) (Timestamp, error) {
	ts, err := types.TimestampFrom(t)
	return Timestamp{ts}, err
}

// ParseTimestamp creates a timestamp from an ISO-3933 string
func ParseTimestamp(s string) (Timestamp, error) {
	ts, err := types.ParseTimestamp(s)
	return Timestamp{ts}, err
}

// Equal implements equality testing for Timestamp
func (vt Timestamp) Equal(rhs Value) bool {
	switch other := rhs.(type) {
	case Timestamp:
		return vt.t.Compare(other.t) == 0
	default:
		return false
	}
}

// Less implements comparison for Timestamp
func (vt Timestamp) Less(rhs Value) (bool, error) {
	switch other := rhs.(type) {
	case Timestamp:
		return vt.t.Compare(other.t) < 0, nil
	default:
		return false, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt Timestamp) IsScalar() bool {
	return true
}

func (vt Timestamp) String() string {
	return vt.t.String()
}

// IsTrue indicates if this Value evaluates to true
func (vt Timestamp) IsTrue() bool {
	return false
}

// T returns the timestamp as a int64 duration in uSec since the start of epoch.
func (vt Timestamp) T() int64 {
	return int64(vt.t)
}
