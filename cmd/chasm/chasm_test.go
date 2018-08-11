package main

import (
	"testing"
)

func TestSimple1(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
		nop
		}
`
	checkParse(t, "Simple1", code, "00 8000 00 88")
}

func TestSimple2(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			nop ; nop instruction
			drop ; drop nothing
		}
`
	checkParse(t, "Simple2", code, "00 8000 0001 88")
}

func TestSimplePush(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			push 0
		}
`
	checkParse(t, "SimplePush", code, "00 8000 20 88")
}

func TestPushB(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			pushb 5 1 2 3 4 5 6 7 8 9 10
			pushb "HI!"
			pushb 0x01 0x02 0x03
			pushb 0x01 0x02 0x03 0x04 0x05 0x06 0x07
			pushb addr(deadbeeffeed5d00d5)
		}
`
	checkParse(t, "PushB", code, `00 8000
	29 0b 050102030405060708090a
	29 03 484921
	29 03 01 02 03
	29 07 01 02 03 04 05 06 07
	29 de ad be ef fe ed 5d 00 d5
	88`)
}

func TestPushT(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			pusht 2018-02-03T01:23:45Z
		}
`
	checkParse(t, "PushT", code, `00 8000
	2c 40 6a e1 72  43 07 02 00
	88`)
}

func TestFunc(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			zero
			call bar 1
		}
		func buzz {
		}
		func bar {
			one
			add
		}
`
	checkParse(t, "Func", code, "00 8000 20 810201 88 8001 88 8002 2a 40 88")
}

func TestSeveralPushes(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			push -1
			push 1
			push 2
			push 12
		}
`
	checkParse(t, "SeveralPushes", code, "00 8000 2b2a2102210c 88")
}

func TestConstants(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			K = 65535
			push K
		}
`
	checkParse(t, "Constants", code, "00 8000 22FFFF 88")
}

func TestUnitaryOpcodes1(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			nop
			drop
			drop2
			dup
			dup2
			swap
			over
			ret
			fail
			zero
			false
			one
			true
			neg1
			now
			rand
			add
			sub
			mul
			div
			mod
			divmod
			muldiv
			not
			neg
			inc
			dec
			index
			len
			append
			extend
			slice
		}
`
	checkParse(t, "Unitary1", code, `
		00 8000
		00 0102 0506 090C
		1011 2020 2a2a 2b2d
		2f40 4142 4344 4546
		4849 4A4B 5051 5253
		5488`)
}

func TestUnitaryOpcodes2(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			choice
			ifz
			ifnz
			else
			endif
			sum
			avg
			max
			min
			pushl
			eq
			gt
			lt
		}
`
	checkParse(t, "Unitary2", code,
		"00 800094 898a8e8f9091929330 4d4e4f 88")
}

func TestBinary(t *testing.T) {
	code := `
		; comment
		context: TEST
		func foo {
			pick 2
			pick 12
			roll 0xA
			field 3
			fieldl 0
			call bar 0
			lookup bar 1
			deco bar 2
			sort 3
			wchoice 1
		}
		func bar {
			nop
		}
`
	checkParse(t, "Binary", code, "00 8000 0D020D0C0E0A 6003 7000 810100 970101 820102 9603 9501 88 8001 00 88")
}

func TestRealistic(t *testing.T) {
	code := `
		; This program pushes a, b, c,
		; and x on the stack and calculates
		; a*x*x + b*x + c
		context: TEST
		func foo {
			A = 3
			B = 5
			C = 7
			X = 21

			push A
			push B
			push C
			push X	; ABCX
			roll 4	; BCXA
			pick 1	; BCXAX
			dup  	; BCXAXX
			mul		; BCXAR
			mul		; BCXR
			roll 4  ; CXRB
			roll 2  ; CRBX
			mul		; CRS
			add		; CR
			add		; R
			ret
		}
`
	checkParse(t, "Realistic", code, `
		00 80 00 21 03 21 05 21 07 21  15 0E 04 0d 01 05 42 42
		0e 04 0e 02 42 40 40 10 88`)
}
