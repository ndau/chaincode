package chain

import (
	"crypto/sha256"
	"math/rand"
)

// SeededRand is a random number generator that conforms to Randomer.
// It accepts a seed at initialization and generates a predictable sequence thereafter.
type SeededRand struct {
	rand *rand.Rand
}

// NewSeededRand returns a new instance of a SeededRand.
func NewSeededRand(ba []byte) (*SeededRand, error) {
	h := sha256.New()
	h.Write(ba)
	seed := int64(0)
	for i, b := range h.Sum(nil) {
		seed ^= int64(b) << uint(i%8)
	}
	return &SeededRand{rand: rand.New(rand.NewSource(seed))}, nil
}

// RandInt implements Randomer for SeededRand.
func (dr *SeededRand) RandInt() (int64, error) {
	r := int64(dr.rand.Uint64())
	return r, nil
}
