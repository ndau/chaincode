package vm

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildVM(t *testing.T, s string) *ChaincodeVM {
	ops := MiniAsm(s)
	assert.Nil(t, ops.IsValid())
	bin := ChasmBinary{"test", "TEST", ops}
	vm, err := New(bin)
	assert.Nil(t, err)
	return vm
}

func buildVMfail(t *testing.T, s string) {
	ops := MiniAsm(s)
	bin := ChasmBinary{"test", "TEST", ops}
	_, err := New(bin)
	assert.NotNil(t, err)
}

// a simple func that can be passed to Run() to print out the steps as it runs
func pr(vm *ChaincodeVM) {
	fmt.Println(vm)
}

func checkStack(t *testing.T, st *Stack, values ...int64) {
	for i := range values {
		n, err := st.PopAsInt64()
		assert.Nil(t, err)
		assert.Equal(t, values[len(values)-i-1], n)
	}
}

func TestMiniAsmBasics(t *testing.T) {
	ops := MiniAsm("neg1 zero one push1 45 push2 01 01 2000-01-01T00:00:00Z")
	bytes := Chaincode{OpNeg1, OpZero, OpOne, OpPush1, 69, OpPush2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	assert.Equal(t, ops, bytes)
}

func TestNop(t *testing.T) {
	vm := buildVM(t, "handler 0 nop enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	assert.Equal(t, vm.Stack().Depth(), 0)
}

func TestPush(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 zero one maxnum minnum push1 45 push2 01 02 ret enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 1, math.MaxInt64, math.MinInt64, 69, 513)
}

// test that simple pushes and sign extension work right
func TestPush1(t *testing.T) {
	vm := buildVM(t, `
	handler 0
	push1 7F ; should be 127
	push1 80 ; should be -128
	push1 FF ; this should be -1
	push1 F0 ; this should be -16
	enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 127, -128, -1, -16)
}

func TestBadHandler(t *testing.T) {
	// this test is designed to catch a handler without a parameter
	buildVMfail(t, `handler`)
}

func TestBadSize(t *testing.T) {
	buildVMfail(t, `
	handler 0
	2D 10
	enddef
	`)
}

// make sure that we check every possible 1-byte value
func TestPush1All(t *testing.T) {
	for i := -128; i < 128; i++ {
		s := fmt.Sprintf("handler 0 push1 %02x enddef", byte(int8(i)))
		vm := buildVM(t, s)
		vm.Init(0)
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), int64(i))
	}
}

// Push the value 5 onto the stack using every possible push operator
func TestPush5AllWays(t *testing.T) {
	vm := buildVM(t, `
	handler 0
	push1 5
	push2 5 0
	push3 5 0 0
	push4 5 0 0 0
	push5 5 0 0 0 0
	push6 5 0 0 0 0 0
	push7 5 0 0 0 0 0 0
	push8 5 0 0 0 0 0 0 0
	enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5, 5, 5, 5, 5, 5, 5, 5)
}

// Push the value -5 onto the stack using every possible push operator
func TestPushMinus5AllWays(t *testing.T) {
	vm := buildVM(t, `
	handler 0
	push1 fb
	push2 fb ff
	push3 fb ff ff
	push4 fb ff ff ff
	push5 fb ff ff ff ff
	push6 fb ff ff ff ff ff
	push7 fb ff ff ff ff ff ff
	push8 fb ff ff ff ff ff ff ff
	enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -5, -5, -5, -5, -5, -5, -5, -5)
}

func TestBigPush(t *testing.T) {
	vm := buildVM(t, `handler 0
		push3 1 2 3
		push4 4 0 0 1
		push5 5 0 0 0 1
		push6 6 0 0 0 0 1
		push7 1 2 3 4 5 6 7
		push8 fb ff ff ff ff ff ff ff enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 197121, 16777220, 4294967301, 1099511627782, 1976943448883713, -5)
}

func TestPushB1(t *testing.T) {
	vm := buildVM(t, "handler 0 pushb 09 41 42 43 44 45 46 47 48 49 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewBytes(nil), v)
	assert.Equal(t, NewBytes([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73}), v)
}

func TestPushB2(t *testing.T) {
	vm := buildVM(t, `handler 0 pushb "ABCDEFGHI" enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewBytes(nil), v)
	assert.Equal(t, NewBytes([]byte{65, 66, 67, 68, 69, 70, 71, 72, 73}), v)
}

func TestDrop(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 7 nop one zero neg1 drop drop2 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 7)
}

func TestDup(t *testing.T) {
	vm := buildVM(t, "handler 0 one push1 2 dup push1 3 dup2 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 2, 2, 3, 2, 3)
}

func TestSwapOverPickRoll(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 swap over pick 4 roll 4 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 3, 2, 3, 0, 1)
}

func TestPickRollEdgeCases(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one pick 0 push1 2 roll 0 push1 3 roll 1 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, 1, 3, 2)
}

func TestTuck(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 tuck 0 tuck 1 tuck 1 tuck 3 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3, 0, 1, 2)
}

func TestTuckFail(t *testing.T) {
	vm := buildVM(t, "handler 0 zero one push1 2 push1 3 tuck 4 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestMath(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 55 dup dup add sub push1 7 push1 6 mul dup push1 3 div dup push1 3 mod enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -85, 42, 14, 2)
}

func TestMul(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 15},
		a{[]int64{5, 5}, 25},
		a{[]int64{3, 5}, 15},
		a{[]int64{12, 4}, 48},
		a{[]int64{5, -3}, -15},
		a{[]int64{5, 0}, 0},
		a{[]int64{0, 5}, 0},
		a{[]int64{-12, -4}, 48},
	}
	vm := buildVM(t, "handler 0 mul enddef")

	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestDiv(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 1},
		a{[]int64{5, 5}, 1},
		a{[]int64{3, 5}, 0},
		a{[]int64{12, 4}, 3},
		a{[]int64{5, -3}, -1},
		a{[]int64{50, 5}, 10},
		a{[]int64{0, 5}, 0},
		a{[]int64{-12, -4}, 3},
	}
	vm := buildVM(t, "handler 0 div enddef")

	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestMod(t *testing.T) {
	type a struct {
		in  []int64
		out int64
	}
	args := []a{
		a{[]int64{5, 3}, 2},
		a{[]int64{5, 5}, 0},
		a{[]int64{3, 5}, 3},
		a{[]int64{12, 4}, 0},
	}
	vm := buildVM(t, "handler 0 mod enddef")
	for a := range args {
		vm.Init(0, NewNumber(args[a].in[0]), NewNumber(args[a].in[1]))
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), args[a].out)
	}
}

func TestDivMod(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 17 push1 7 divmod enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2, 3)
}

func TestMulDiv(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 64 push1 11 push1 19 muldiv enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 68)
}

func TestMulDivBig(t *testing.T) {
	vm := buildVM(t, "handler 0 push8 00 00 b2 d3 59 5b f0 06 push6 00 00 00 00 00 01 push6 00 00 00 00 00 02 muldiv enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 250000000000000000)
}

func TestMathErrors(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 55 zero div enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err, "divide by zero")

	vm = buildVM(t, "handler 0 push1 55 zero mod enddef")
	vm.Init(0)
	err = vm.Run(nil)
	assert.NotNil(t, err, "mod by zero")

	vm = buildVM(t, "handler 0 push1 55 zero divmod enddef")
	vm.Init(0)
	err = vm.Run(nil)
	assert.NotNil(t, err, "divmod by zero")

	vm = buildVM(t, "handler 0 push1 55 push1 2 zero muldiv enddef")
	vm.Init(0)
	err = vm.Run(nil)
	assert.NotNil(t, err, "muldiv by zero")
}

func TestMathOverflows(t *testing.T) {
	vm := buildVM(t, "handler 0 push8 7a bb cc dd ee ff 99 88 push1 7f mul enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err, "mul overflow")
	vm = buildVM(t, "handler 0 push8 7f bb cc dd ee ff 99 88 push8 7f bb cc dd ee ff 99 88 add enddef")
	vm.Init(0)
	err = vm.Run(nil)
	assert.NotNil(t, err, "add overflow")
	vm = buildVM(t, "handler 0 push8 7f bb cc dd ee ff 99 78 push8 ff bb cc dd ee ff 99 88 sub enddef")
	vm.Init(0)
	err = vm.Run(nil)
	assert.NotNil(t, err, "sub overflow")
}

func TestNot(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 7 not zero not pushl not enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1, -1)
}

func TestNotNegIncDec(t *testing.T) {
	vm := buildVM(t, "handler 0 push1 7 not dup not push1 8 neg push1 4 inc push1 6 dec enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1, -8, 5, 5)
}

func TestIf1(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf2(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifnz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf3(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf4(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 else push1 42 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf5(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf6(t *testing.T) {
	vm := buildVM(t, "handler 0 zero ifnz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf7(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf8(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIfNested1(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 zero ifz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 21, 17)
}

func TestIfNested2(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifz push1 13 zero ifz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIfNested3(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz push1 13 zero ifnz push1 15 else push1 13 endif endif push1 11 enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 19, 17)
}

func TestIfNested4(t *testing.T) {
	vm := buildVM(t, `
	handler 0 					; X
		dup						; X X
		push1 4					; X X 4
		lt						; X (X < 4)
		ifnz 					; X
			dup					; X X
			push1 2				; X X 2
			lt					; X (X < 2)
			ifnz 				; X
				dup				; X X
				push1 1			; X X 1
				lt				; X (X < 1)
				ifnz 			; X
					push1 40	; X 40
				else
					push1 41	; X 41
				endif
			else
				dup				; X X
				push1 3			; X X 3
				lt				; X (X < 3)
				ifnz 			; X
					push1 42	; X 42
				else
					push1 43	; X 43
				endif
			endif
		else
			dup					; X X
			push1 6				; X X 6
			lt					; X (X < 6)
			ifnz 				; X
				dup				; X X
				push1 5			; X X 5
				lt				; X (X < 5)
				ifnz 			; X
					push1 44	; X 44
				else
					push1 45	; X 45
				endif
			else
				dup				; X X
				push1 7			; X X 7
				lt				; X (X < 7)
				ifnz 			; X
					push1 46	; X 46
				else
					push1 47	; X 47
				endif
			endif
		endif
	enddef
	`)
	for i := int64(0); i < 8; i++ {
		vm.Init(0, NewNumber(i))
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), i, 0x40+i)
	}
}

func buildNestedConditional(mask, min, max int) string {
	code := ""
	tmplN := `dup						; X X
		push1 %[1]x					; X X %[1]d
		and						; X (x & %[1]d)
		ifz 					; X
		`
	tmpl1 := `dup					; X X
		push1 1					; X X 1
		and						; X (x & 1)
		ifz					; X
			push1 %x
		else
			push1 %x
		endif
		`
	if mask == 1 {
		code = fmt.Sprintf(tmpl1, min, min+1)
	} else {
		code = fmt.Sprintf(tmplN, mask)
		code += buildNestedConditional(mask>>1, min, (min+max)/2)
		code += "else\n"
		code += buildNestedConditional(mask>>1, (min+max)/2, max)
		code += "endif\n"
	}
	return code
}

// formatLine simply formats a single line of asm
// so that we have a hope of reading/debugging the generated code
// It's relatively stupid and will definitely fail if your line
// contains a semicolon that is not the start of a comment
func formatLine(s string, indent, commentindent int) (string, int) {
	const indentstep = 2
	// this doesn't work with quoted strings that contain ;
	p := regexp.MustCompile("[ \t]*([A-Za-z0-9_]+)?([^;]*)?(;.*)?")
	parts := p.FindStringSubmatch(s)
	// fmt.Printf("%#v\n", parts)
	parts[2] = strings.TrimSpace(parts[2])
	if parts[0] == "" && parts[2] != "" {
		return fmt.Sprintf("%*s%s", indent, "", parts[3]), indent
	}
	newindent := indent
	switch strings.ToLower(parts[1]) {
	case "handler", "ifz", "ifnz", "def":
		newindent += indentstep
	case "else":
		indent -= indentstep
	case "enddef", "endif":
		indent -= indentstep
		newindent -= indentstep
	}
	code := fmt.Sprintf("%*s%s %s", indent, "", parts[1], parts[2])
	return fmt.Sprintf("%-*s%s", commentindent, code, parts[3]), newindent
}

// formatChaincode formats a block of chaincode
func formatChaincode(s string) string {
	lines := strings.Split(s, "\n")
	newlines := make([]string, len(lines))
	commentind := 40
	ind := 0
	for i := range lines {
		newlines[i], ind = formatLine(lines[i], ind, commentind)
	}
	return strings.Join(newlines, "\n")
}

// This test generates a deeply nested structure that
// simply does a binary test of which bits are set in a value;
// we can run it up to 4 levels deep without running out of
// bytes in the VM. It evaluates that the if..else...endif structure
// works properly even with deep nesting.
func TestIfNestedDeep(t *testing.T) {
	const depth = 4
	mask := 1 << depth
	s := "handler 0\n"
	s += buildNestedConditional(mask>>1, 0, mask)
	s += "enddef\n"
	// fmt.Println(formatChaincode(s))

	vm := buildVM(t, s)
	for i := int64(0); i < int64(mask); i++ {
		vm.Init(0, NewNumber(i))
		err := vm.Run(nil)
		assert.Nil(t, err)
		checkStack(t, vm.Stack(), i, i)
	}

}

func TestIfNull1(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz endif enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack())
}

func TestIfNull2(t *testing.T) {
	vm := buildVM(t, "handler 0 one ifnz else endif enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack())
}

func TestCompares1(t *testing.T) {
	vm := buildVM(t, "handler 0 one neg1 eq one neg1 lt one neg1 gt one neg1 lte one neg1 gte enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, -1, 0, -1)
}

func TestCompares2(t *testing.T) {
	vm := buildVM(t, "handler 0 one one eq one one lt one one gt one one lte one one gte enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 0, -1, -1)
}

func TestCompares3(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 one eq neg1 one lt neg1 one gt neg1 one lte neg1 one gte enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1, 0, -1, 0)
}

func TestCompares4(t *testing.T) {
	vm := buildVM(t, "handler 0 neg1 pushb 8 1 2 3 4 5 6 7 8 eq enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestCompares5(t *testing.T) {
	vm := buildVM(t, `handler 0 pushb "hello" pushb "hi" dup2 eq pick 2 pick 2 lt enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1)
}

func TestCompareLists1(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup dup eq swap dup dup gt swap dup gte enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, -1)
}

func TestCompareLists2(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup one append dup pick 2 eq swap roll 2 gt enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1)
}

func TestCompareLists3(t *testing.T) {
	vm := buildVM(t, `handler 0 pushl zero append one append dup one append swap dup pick 2 eq swap roll 2 gt enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0)
}

func TestCompare7(t *testing.T) {
	vm := buildVM(t, "handler 0 dup zero index pick 1 one index eq enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewTestStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)
}

func TestCompareTimestampGt(t *testing.T) {
	// This checks that timestamps are correctly ordered by gt
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		gt
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		gt
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		gt
		enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 0)
}

func TestCompareTimestampLt(t *testing.T) {
	// This checks that timestamps are correctly ordered by lt
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		lt
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		lt
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		lt
		enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, -1, 0)
}

func TestCompareTimestampEq(t *testing.T) {
	// This checks that timestamps are correctly ordered by eq
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z pusht 2018-01-01T00:00:00Z
		eq
		pusht 2018-01-01T00:00:00Z pusht 2018-07-18T00:00:00Z
		eq
		pusht 2018-07-18T00:00:00Z pusht 2018-07-18T00:00:00Z
		eq
		enddef`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, -1)
}

func TestTimestampNegativePush(t *testing.T) {
	// This checks that a timestamp built from bytes cannot be negative
	vm := buildVM(t, `
		handler 0
		pusht 10 20 30 40 50 60 70 80
		enddef
		`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestTimestamp1(t *testing.T) {
	// This checks that subtracting timestamps returns the appropriate
	// positive value
	vm := buildVM(t, `
		handler 0
		pusht 2018-07-18T00:00:00Z
		pusht 2018-01-01T00:00:00Z
		sub
		push3 40 42 0f
		div
		push1 3C
		dup
		mul
		push1 18
		mul
		div
		enddef
		`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 198)
}

func TestTimestampNegativeSub(t *testing.T) {
	// This checks that subtracting timestamps can return the appropriate
	// negative value.
	vm := buildVM(t, `
		handler 0
		pusht 2018-01-01T00:00:00Z
		pusht 2018-07-18T00:00:00Z
		sub
		push3 40 42 0f
		div
		push1 3C
		dup
		mul
		push1 18
		mul
		div
		enddef
		`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -198)
}

func TestTimestampInjectedNow(t *testing.T) {
	vm := buildVM(t, `
		handler 0
		now
		pusht 2018-01-01T00:00:00Z
		sub
		enddef
		`)
	ts, err := ParseTimestamp("2018-01-02T03:04:05Z")
	assert.Nil(t, err)
	now, err := NewCachingNow(ts)
	assert.Nil(t, err)
	vm.SetNow(now)
	vm.Init(0)
	err = vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 97445000000)
}

func TestTimestampDefaultNow(t *testing.T) {
	// This checks that the default now operation returns a date stamp
	// between 1/1/18 and 2/2/22 (which will fail someday but not for a
	// few years)
	vm := buildVM(t, `
		handler 0
		now
		dup
		pusht 2018-01-01T00:00:00Z
		lt
		swap
		now
		sub
		zero
		eq
		pusht 2022-02-02T22:22:22Z
		now
		gt
		enddef
		`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, -1)
}

func TestInjectedRand(t *testing.T) {
	vm := buildVM(t, "handler 0 rand rand eq rand rand eq enddef")
	r, err := NewCachingRand()
	assert.Nil(t, err)
	vm.SetRand(r)
	vm.Init(0)
	err = vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, -1)
}

func TestDefaultRand(t *testing.T) {
	vm := buildVM(t, "handler 0 rand rand eq rand rand eq enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0)
}

func TestList1(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl push1 0d append push1 7 append dup len swap dup one index swap push1 2 neg index enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2, 7, 13)
}

func TestExtend(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 7 append dup zero append swap extend dup len swap push1 2 index enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5, 0)
}

func TestSlice(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append one push1 3 slice len enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestSlice2(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append dup len one sub zero swap slice len enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestSum(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append push1 3 append sum enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 6)
}

type seededRand struct {
	n int64
}

// RandInt implements Randomer for seededRand
func (r seededRand) RandInt() (int64, error) {
	return r.n, nil
}

func TestChoice(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl zero append one append push1 2 append push1 3 append choice enddef")
	r := seededRand{n: 12345}
	vm.SetRand(r)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestWChoice1(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	for i := int64(0); i < 6; i++ {
		st := NewTestStruct(NewNumber(i))
		l = l.Append(st)
	}
	vm.Init(0, l)

	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 4)
}

func TestWChoice2(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	for i := int64(0); i < 6; i++ {
		st := NewTestStruct(NewNumber(6 - i))
		l = l.Append(st)
	}
	vm.Init(0, l)

	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5)
}

func TestWChoiceErr(t *testing.T) {
	vm := buildVM(t, "handler 0 wchoice 0 field 0 enddef")
	r := seededRand{n: math.MaxInt64 / 2}
	vm.SetRand(r)

	l := NewList()
	vm.Init(0, l)

	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestAvg(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 7 append push1 16 append avg enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 10)
}

func TestMin(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 2 append push1 3 append min enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1)
}

func TestMax(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl one append push1 2 append push1 3 append max enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestAvgFail(t *testing.T) {
	vm := buildVM(t, "handler 0 pushl avg enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestField(t *testing.T) {
	vm := buildVM(t, "handler 0 field 2 enddef")
	st := NewTestStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(0, st)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 27)
}

func TestIsField(t *testing.T) {
	vm := buildVM(t, "handler 0 dup isfield 2 swap isfield 3 enddef")
	st := NewTestStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(0, st)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0)
}

func TestFieldFail(t *testing.T) {
	vm := buildVM(t, "handler 0 field 9 enddef")
	st := NewTestStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(0, st)
	err := vm.Run(nil)
	assert.NotNil(t, err)

	vm = buildVM(t, "handler 0 isfield 9 enddef")
	vm.Init(0, NewNumber(27))
	err = vm.Run(nil)
	assert.NotNil(t, err)
}

func TestFieldL(t *testing.T) {
	vm := buildVM(t, "handler 0 fieldl 2 one index enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewTestStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5)
}

func TestFieldLFail(t *testing.T) {
	vm := buildVM(t, "handler 0 fieldl 9 one index enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewTestStruct(NewNumber(3*i), NewNumber(3*i+1), NewNumber(3*i+2))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestSortFields(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 2 push1 3 index field 1 enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewTestStruct(NewNumber(2*i), NewNumber(3*i+1), NewNumber(4*(6-i)))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 4)
}

func TestSortFail(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 6 push1 3 index field 1 enddef")
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewTestStruct(NewNumber(2*i), NewNumber(3*i+1), NewNumber(4*(6-i)))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestNestingFail1(t *testing.T) {
	buildVMfail(t, "def 1 nop enddef")
	buildVMfail(t, "handler 0 nop enddef handler 0 nop enddef")
	buildVMfail(t, "handler 0 nop enddef def 2 nop enddef")
	buildVMfail(t, "handler 0 ifz enddef")
	buildVMfail(t, "handler 0 ifnz enddef")
	buildVMfail(t, "handler 0 enddef enddef")
	buildVMfail(t, "handler 0 ifz else else enddef enddef")
	buildVMfail(t, "handler 0 push8 enddef")
}

func TestCall1(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 enddef def 0 1 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestCall2(t *testing.T) {
	vm := buildVM(t, `
		handler 0 one call 0 enddef
		def 0 1 push1 2 call 1 enddef
		def 1 2 add enddef
	`)
	vm.Init(0)
	err := vm.Run(nil)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestCallFail1(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 1 enddef def 0 1 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestCallFail2(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 enddef def 0 2 push1 2 add enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestCallFail3(t *testing.T) {
	vm := buildVM(t, "handler 0 one call 0 enddef def 0 1 push1 2 add drop enddef")
	vm.Init(0)
	err := vm.Run(nil)
	assert.NotNil(t, err)
}

func TestSizeFail1(t *testing.T) {
	s := `handler 0
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
		pushb 20 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4
	enddef`
	// this VM has big data but small ish real code
	// so by default it should load
	vm := buildVM(t, s)
	assert.NotNil(t, vm, "no error")
	// but if we put the maximum code size to 30
	// it should fail with code too long
	SetMaxLengths(30, 1024)
	buildVMfail(t, s)
	// and if we put max data size to 256 it should fail with data too long
	SetMaxLengths(256, 256)
	buildVMfail(t, s)
}
