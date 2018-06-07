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

// NewTimestamp creates a timestamp from an int64 representation of one
func NewTimestamp(n int64) Timestamp {
	return Timestamp{types.Timestamp(n)}
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

// Compare implements comparison for Timestamp
func (vt Timestamp) Compare(rhs Value) (int, error) {
	switch other := rhs.(type) {
	case Timestamp:
		return vt.t.Compare(other.t), nil
	default:
		return 0, ValueError{"comparing incompatible types"}
	}
}

// IsScalar indicates if this Value is a scalar value type
func (vt Timestamp) IsScalar() bool {
	return true
}

func (vt Timestamp) String() string {
	return vt.t.String()
}

// T returns the timestamp as a int64 duration in uSec since the start of epoch.
func (vt Timestamp) T() int64 {
	return int64(vt.t)
}
