package vm

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	cryptorand "crypto/rand"
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
	b, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
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
	b, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
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
	return NewTimestampFromTime(time.Now())
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

// FractionLess compares two ratios of int64 values (fractions) without exceeding an int64.
func FractionLess(n1, d1, n2, d2 int64) bool {
	return compareRatios(n1, d1, n2, d2, true)
}

// compareRatios compares two ratios of int64 values (fractions) without exceeding an int64.
// It does so by computing the continued fraction representation of each fraction
// simultaneously, but stops as soon as they differ.
// The algorithm was found here:
// https://janmr.com/blog/2014/05/comparing-rational-numbers-without-overflow/
func compareRatios(n1, d1, n2, d2 int64, less bool) bool {
	i1 := n1 / d1
	r1 := n1 % d1
	i2 := n2 / d2
	r2 := n2 % d2
	switch {
	case n2 == 0:
		return false
	case n1 == 0:
		return true
	case less && i1 < i2,
		!less && i1 > i2:
		return true
	case less && i1 > i2,
		!less && i1 < i2:
		return false
	case r1 == 0:
		return false
	case r2 == 0:
		return true
	default:
		return compareRatios(d1, r1, d2, r2, !less)
	}
}
