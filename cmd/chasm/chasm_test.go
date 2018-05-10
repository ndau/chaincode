package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkParse(t *testing.T, name string, code string, result string) {
	sn, err := Parse(name, []byte(code))
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	b := sn.(Script).bytes()
	bcheck(t, b, result)
}

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
	checkParse(t, "Unitary1", code, "000001020506090D101120202a2a2b2d2f4041424344454647485051525354")
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
		}
`
	checkParse(t, "Unitary2", code, "006070949596978081828890919293")
}
