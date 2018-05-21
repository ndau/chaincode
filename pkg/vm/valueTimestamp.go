package vm

import (
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

// NewTimestampFromTime returns a timestamp taken from a time.Time struct in Go.
func NewTimestampFromTime(t time.Time) Timestamp {
	epoch, err := time.Parse(TimestampFormat, EpochStart)
	if err != nil {
		panic("Epoch isn't a valid timestamp!")
	}
	uSec := uint64(t.Sub(epoch).Nanoseconds() / 1000)
	return NewTimestamp(uSec)
}

// ParseTimestamp creates a timestamp from an ISO-3933 string
func ParseTimestamp(s string) (Timestamp, error) {
	ts, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return NewTimestamp(0), err
	}
	return NewTimestampFromTime(ts), nil
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
	epoch, _ := time.Parse(TimestampFormat, EpochStart)
	t := epoch.Add(time.Duration(vt.t) * time.Microsecond)
	return t.Format(TimestampFormat)
}

// T returns the timestamp as a uint64 duration in uSec since the start of epoch.
func (vt Timestamp) T() uint64 {
	return vt.t
}
