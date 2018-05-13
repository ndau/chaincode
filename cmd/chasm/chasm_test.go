package main

import (
	"testing"
)

func TestSimple1(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
		nop
		}
`
	checkParse(t, "Simple1", code, "0000")
}

func TestSimple2(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			nop ; nop instruction
			drop ; drop nothing
		}
`
	checkParse(t, "Simple2", code, "000001")
}

func TestSimplePush(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			push 0
		}
`
	checkParse(t, "SimplePush", code, "0020")
}

func TestSeveralPushes(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			push -1
			push 1
			push 2
			push 12
		}
`
	checkParse(t, "SeveralPushes", code, "002b2a2102210c")
}

func TestConstants(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			K = 65535
			push K
		}
`
	checkParse(t, "Constants", code, "0022FFFF")
}

func TestUnitaryOpcodes1(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
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
		0000 0102 0506 090D
		1011 2020 2a2a 2b2d
		2f40 4142 4344 4546
		4748 5051 5253 54`)
}

func TestUnitaryOpcodes2(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			field
			fieldl
			choice
			wchoice
			sort
			lookup
			ifz
			ifnz
			else
			end
			sum
			avg
			max
			min
			pushl
		}
`
	checkParse(t, "Unitary2", code,
		"00607094959697808187889091929330")
}

func TestBinary(t *testing.T) {
	code := `
		; comment
		context: TEST
		{
			pick 2
			pick 12
			roll 0xA
		}
`
	checkParse(t, "Binary", code, "000E020E0C0F0A")
}

func TestRealistic(t *testing.T) {
	code := `
		; This program pushes a, b, c,
		; and x on the stack and calculates
		; a*x*x + b*x + c
		context: TEST
		{
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
		00 21 03 21 05 21 07 21  15 0f 04 0e 01 05 42 42
		0f 04 0f 02 42 40 40 10`)
}
