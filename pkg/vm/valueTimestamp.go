package vm

import (
	"strconv"
	"time"
)

// Constants for time
const (
	EpochStart      = "2018-01-01T00:00:00Z"
	TimestampFormat = "2006-01-02T15:04:05Z"
)

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

// IsScalar indicates if this Value is a scalar value type
func (vt Timestamp) IsScalar() bool {
	return true
}

func (vt Timestamp) String() string {
	return strconv.FormatUint(vt.t, 16)
}

// T returns the timestamp as a uint64 duration in uSec since the start of epoch.
func (vt Timestamp) T() uint64 {
	return vt.t
}
