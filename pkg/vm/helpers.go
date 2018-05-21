package vm

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// Some masking values
const (
	ByteMask byte = 0xFF
	HighBit  byte = 0x80
)

// ToBytesU returns an array of 8 bytes encoding n as a uint in little-endian form
func ToBytesU(n uint64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// ToBytes returns an array of 8 bytes encoding n as a signed value in little-endian form
func ToBytes(n int64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}

// DefaultRand is the default random number generator, which picks a random number from 0 to MaxInt64
// using crypto/rand.
type DefaultRand struct{}

// NewDefaultRand returns a new instance of a DefaultRand
func NewDefaultRand() (*DefaultRand, error) {
	return &DefaultRand{}, nil
}

// RandInt implements Randomer for DefaultRand
func (dr *DefaultRand) RandInt() (int64, error) {
	b, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return -1, err
	}
	r := b.Int64()
	return r, nil
}

// CachingRand is a random number generator that picks a random number once and then
// returns the same value over and over
type CachingRand struct {
	r int64
}

// NewCachingRand returns an initialized instance of a CachingRand
func NewCachingRand() (*CachingRand, error) {
	b, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	cr := CachingRand{r: b.Int64()}
	return &cr, nil
}

// RandInt implements Randomer for CachingRand
func (dr *CachingRand) RandInt() (int64, error) {
	return dr.r, nil
}

// DefaultNow is the default time object, which returns a Timestamp based on the current time.
type DefaultNow struct{}

// NewDefaultNow returns a new instance of a DefaultNow
func NewDefaultNow() (*DefaultNow, error) {
	return &DefaultNow{}, nil
}

// Now implements Nower for DefaultNow
func (dr *DefaultNow) Now() (Timestamp, error) {
	epoch, err := time.Parse(TimestampFormat, EpochStart)
	if err != nil {
		return NewTimestamp(0), err
	}
	t := time.Now()
	uSec := uint64(t.Sub(epoch).Nanoseconds() / 1000)
	return NewTimestamp(uSec), nil
}

// CachingNow is a Nower that
// returns the same value over and over
type CachingNow struct {
	t Timestamp
}

// NewCachingNow returns an initialized instance of a CachingNow
// given a Timestamp it should use
func NewCachingNow(ts Timestamp) (*CachingNow, error) {
	cn := CachingNow{t: ts}
	return &cn, nil
}

// Now implements Nower for CachingNow
func (cn *CachingNow) Now() (Timestamp, error) {
	return cn.t, nil
}
