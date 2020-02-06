package chain

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"crypto/sha256"
	"math/rand"

	"github.com/bszcz/mt19937_64" // Use standard package rather than our old fork
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
	source := mt19937_64.New()
	source.Seed(seed)
	return &SeededRand{rand: rand.New(source)}, nil
}

// Seed sets the seed for the current SeededRand object. Note that seeding is expensive;
// doing NewSeededRand takes only a tiny amount longer than reseeding an existing random
// number generator.
func (sr *SeededRand) Seed(seed int64) {
	sr.rand.Seed(seed)
}

// RandInt implements Randomer for SeededRand.
func (sr *SeededRand) RandInt() (int64, error) {
	return sr.rand.Int63(), nil
}
