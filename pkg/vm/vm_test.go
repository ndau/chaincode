package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildVM(t *testing.T, s string) *ChaincodeVM {
	ops := miniAsm(s)
	bin := ChasmBinary{"test", "", "TEST", ops}
	vm, err := New(bin)
	assert.Nil(t, err)
	return vm
}

func checkStack(t *testing.T, st *Stack, values ...int64) {
	for i := range values {
		n, err := st.PopAsInt64()
		assert.Nil(t, err)
		assert.Equal(t, values[len(values)-i-1], n)
	}
}

func TestMiniAsmBasics(t *testing.T) {
	ops := miniAsm("neg1 zero one push1 45 push2 01 01 2018-01-01T00:00:00Z")
	bytes := []Opcode{0, OpNeg1, OpZero, OpOne, OpPush1, 69, OpPush2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	assert.Equal(t, ops, bytes)
}

func TestNop(t *testing.T) {
	vm := buildVM(t, "nop")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	assert.Equal(t, vm.Stack().Depth(), 0)
}

func TestPush(t *testing.T) {
	vm := buildVM(t, "neg1 zero one push1 45 push2 01 02 ret")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 1, 69, 513)
}

func TestBigPush(t *testing.T) {
	vm := buildVM(t, "push3 1 2 3 push7 1 2 3 4 5 6 7 push8 fb ff ff ff ff ff ff ff")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 197121, 1976943448883713, -5)
}

func TestPush64(t *testing.T) {
	vm := buildVM(t, "push64 1 2 3 4 5 6 7 8")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewID(0), v)
	assert.Equal(t, NewID(578437695752307201), v)
}

func TestDrop(t *testing.T) {
	vm := buildVM(t, "push1 7 nop one zero neg1 drop drop2")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 7)
}

func TestDup(t *testing.T) {
	vm := buildVM(t, "one push1 2 dup push1 3 dup2")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 2, 2, 3, 2, 3)
}

func TestSwapOverPickRoll(t *testing.T) {
	vm := buildVM(t, "zero one push1 2 push1 3 swap over pick 4 roll 4")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 3, 2, 3, 0, 1)
}

func TestMath(t *testing.T) {
	vm := buildVM(t, "push1 55 dup dup add sub push1 7 push1 6 mul dup push1 3 div dup push1 3 mod")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -85, 42, 14, 2)
}

func TestMathErrors(t *testing.T) {
	vm := buildVM(t, "push1 55 zero div")
	vm.Init()
	err := vm.Run(false)
	assert.NotNil(t, err)
	vm = buildVM(t, "push1 55 zero mod")
	vm.Init()
	err = vm.Run(false)
	assert.NotNil(t, err)
}

func TestNotNegIncDec(t *testing.T) {
	vm := buildVM(t, "push1 7 not dup not push1 8 neg push1 4 inc push1 6 dec")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, -8, 5, 5)
}

func TestIf1(t *testing.T) {
	vm := buildVM(t, "zero ifz push1 13 else push1 42 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf2(t *testing.T) {
	vm := buildVM(t, "zero ifnz push1 13 else push1 42 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf3(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 else push1 42 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf4(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 else push1 42 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf5(t *testing.T) {
	vm := buildVM(t, "zero ifz push1 13 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf6(t *testing.T) {
	vm := buildVM(t, "zero ifnz push1 13 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf7(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf8(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIfNested1(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 zero ifz push1 15 else push1 13 end end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 21, 17)
}

func TestIfNested2(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 zero ifz push1 15 else push1 13 end end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIfNested3(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 zero ifnz push1 15 else push1 13 end end push1 11")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 19, 17)
}

func TestCompares1(t *testing.T) {
	vm := buildVM(t, "one neg1 eq one neg1 lt one neg1 gt")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, 0)
}

func TestCompares2(t *testing.T) {
	vm := buildVM(t, "one one eq one one lt one one gt")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 0, 0)
}

func TestCompares3(t *testing.T) {
	vm := buildVM(t, "neg1 one eq neg1 one lt neg1 one gt")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 0, 1)
}

func TestCompares4(t *testing.T) {
	vm := buildVM(t, "neg1 push64 1 2 3 4 5 6 7 8 eq")
	vm.Init()
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestTimestamp(t *testing.T) {
	vm := buildVM(t, `
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
		`)
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 198)
}

func TestList1(t *testing.T) {
	vm := buildVM(t, "pushl one append push1 7 append dup len swap one index")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2, 7)
}

func TestExtend(t *testing.T) {
	vm := buildVM(t, "pushl one append push1 7 append dup zero append swap extend dup len swap push1 2 index")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 5, 0)
}

func TestSlice(t *testing.T) {
	vm := buildVM(t, "pushl zero append one append push1 2 append one push1 3 slice len")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestSum(t *testing.T) {
	vm := buildVM(t, "pushl zero append one append push1 2 append push1 3 append sum")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 6)
}

func TestAvg(t *testing.T) {
	vm := buildVM(t, "pushl one append push1 7 append push1 16 append avg")
	vm.Init()
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 10)
}

func TestAvgFail(t *testing.T) {
	vm := buildVM(t, "pushl avg")
	vm.Init()
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestField(t *testing.T) {
	vm := buildVM(t, "field 2")
	st := NewStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(st)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 27)
}

func TestFieldFail(t *testing.T) {
	vm := buildVM(t, "field 9")
	st := NewStruct(NewNumber(3), NewNumber(9), NewNumber(27))
	vm.Init(st)
	err := vm.Run(false)
	assert.NotNil(t, err)
}
