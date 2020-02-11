package chain

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestSeededRand(t *testing.T) {
	r, err := NewSeededRand([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Unable to create SeededRand")
	}
	allbits := int64(0)
	const nsamples = 20
	samples := make([]int64, nsamples)
	for i := 0; i < nsamples; i++ {
		v, err := r.RandInt()
		if err != nil {
			t.Errorf("Seeded rand returned error %s", err)
		}
		if v < 0 {
			t.Errorf("negative value %d returned from SeededRand", v)
		}
		allbits |= v
		samples[i] = v
	}
	// A couple of really basic "is it doing anything blatantly stupid" tests:
	// this would catch something like returning Int32s by accident
	if allbits != math.MaxInt64 {
		t.Errorf("Not all bits were set by SeededRand, %016x", allbits)
	}
	// after a decent number of tries, the average should be pretty close to MaxInt64/2
	sum := big.NewInt(0)
	for _, s := range samples {
		sum.Add(sum, big.NewInt(s))
	}
	mean := sum.Div(sum, big.NewInt(nsamples)).Int64()
	expectedMean := int64(math.MaxInt64 / 2)
	deviationFromExpectedMean := expectedMean - mean
	if deviationFromExpectedMean < 0 {
		deviationFromExpectedMean *= -1
	}
	if deviationFromExpectedMean > expectedMean/10 {
		t.Errorf("deviation from expected mean after %d samples was more than 10%% -- %d/%d",
			nsamples, deviationFromExpectedMean, expectedMean)
		fmt.Println(deviationFromExpectedMean, math.MaxInt64/(2*10))
	}
}

func TestThatSeedsMatter(t *testing.T) {
	r1, err := NewSeededRand([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Unable to create SeededRand")
	}
	r2, err := NewSeededRand([]byte{1, 2, 3, 5})
	if err != nil {
		t.Fatalf("Unable to create SeededRand")
	}
	n1, err := r1.RandInt()
	if err != nil {
		t.Fatalf("error on r1")
	}
	n2, err := r2.RandInt()
	if err != nil {
		t.Fatalf("error on r2")
	}
	if n1 == n2 {
		t.Fatalf("seed didn't create different results!: %d %d ", n1, n2)
	}
}

func TestSameSeeds(t *testing.T) {
	r1, err := NewSeededRand([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Unable to create SeededRand")
	}
	r2, err := NewSeededRand([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatalf("Unable to create SeededRand")
	}
	n1, err := r1.RandInt()
	if err != nil {
		t.Fatalf("error on r1")
	}
	n2, err := r2.RandInt()
	if err != nil {
		t.Fatalf("error on r2")
	}
	if n1 != n2 {
		t.Fatalf("same seeds created different results!: %d %d ", n1, n2)
	}
}

var x int64
var r1 *SeededRand

// On a test machine this benchmarked out to less than 100 nanoSec per generated random number
func BenchmarkRandoms(b *testing.B) {
	r1, _ := NewSeededRand([]byte{1, 2, 3, 4})

	for i := 0; i < b.N; i++ {
		x, _ = r1.RandInt()
	}
}

// However, if we have to set up every time, it's 2.2 uSec per random number
func BenchmarkNewRandoms(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1, _ = NewSeededRand([]byte{1, 2, 3, 4})
		x, _ = r1.RandInt()
	}
}

// This is slightly faster than setting up NewSeededRand -- 1.3 uSec each.
func BenchmarkNewRandomSeed(b *testing.B) {
	r1, _ = NewSeededRand([]byte{1, 2, 3, 4})
	for i := 0; i < b.N; i++ {
		r1.Seed(1234567)
		x, _ = r1.RandInt()
	}
}
