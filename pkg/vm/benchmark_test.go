package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

import (
	"strings"
	"testing"
)

var d int
var r int64

func checkd(expected int, b *testing.B) {
	if d != expected {
		b.Errorf("Depth should have been %d, was %d", expected, d)
	}
}

func checkr(expected int64, b *testing.B) {
	if r != expected {
		b.Errorf("Result should have been %d, was %d", expected, r)
	}
}

func benchmarkVM(s string, b *testing.B) {
	ops := MiniAsm(s)
	bin := ChasmBinary{"test", "TEST", ops}
	benchmarkBin(bin, b)
}

func benchmarkBin(bin ChasmBinary, b *testing.B) {
	vm, err := New(bin)
	if err != nil {
		b.Errorf("New() had an error: %s", err)
		return
	}
	// the setup above can be expensive, so make sure we're only benchmarking the runtime
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vm.Init(0)
		err := vm.Run(nil)
		if err != nil {
			b.Errorf("Run() had an error:%s", err)
			return
		}
		d = vm.Stack().Depth()
		if d > 0 {
			r, _ = vm.Stack().PopAsInt64()
		}
	}
}

func benchmarkN(n int, instrs string, b *testing.B) {
	prg := "handler 0 " + strings.Repeat(instrs, n) + "enddef"
	benchmarkVM(prg, b)
}

func BenchmarkNop1(b *testing.B) {
	benchmarkN(1, "nop ", b)
}

func BenchmarkNop10(b *testing.B) {
	benchmarkN(10, "nop ", b)
}

func BenchmarkNop100(b *testing.B) {
	benchmarkN(100, "nop ", b)
}

func BenchmarkAdd1(b *testing.B) {
	benchmarkN(1, "one one add ", b)
}

func BenchmarkAdd10(b *testing.B) {
	benchmarkN(10, "one one add ", b)
	checkd(10, b)
	checkr(2, b)
}

func BenchmarkBigAdd(b *testing.B) {
	prog := `
		push3 1 2 3
		push4 4 0 0 1
		push5 5 0 0 0 1
		push6 6 0 0 0 0 1
		push7 1 2 3 4 5 6 7
		push8 fb ff ff ff ff ff ff ff
		add
		add
		add
		add
		add
`
	benchmarkN(1, prog, b)
	checkd(1, b)
	// this is the sum of all those values: 197121 + 16777220 + 4294967301 + 1099511627782 + 1976943448883713 - 5
	checkr(1978047272453132, b)
}

func BenchmarkMul1(b *testing.B) {
	benchmarkN(1, "one one mul ", b)
}

func BenchmarkMul10(b *testing.B) {
	benchmarkN(10, "one one mul ", b)
}

func BenchmarkQuad(b *testing.B) {
	prog := strings.NewReader(`{
		"name": "examples/quadratic.chasm",
		"comment": "Test of quadratic",
		"data": "oAAhAyEFIQchFQ4DDQEFQkIOAw4CQkBAiA=="
	} `)
	bin, err := Deserialize(prog)
	if err != nil {
		b.Fatal("Unable to parse quad prog")
		return
	}
	benchmarkBin(bin, b)
}
