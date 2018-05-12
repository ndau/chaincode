package vm

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var opcodeMap = map[string]Opcode{
	"nop":     OpNop,
	"drop":    OpDrop,
	"drop2":   OpDrop2,
	"dup":     OpDup,
	"dup2":    OpDup2,
	"swap":    OpSwap,
	"over":    OpOver,
	"pick":    OpPick,
	"roll":    OpRoll,
	"ret":     OpRet,
	"fail":    OpFail,
	"zero":    OpZero,
	"false":   OpFalse,
	"pushN":   OpPushN,
	"push1":   OpPush1,
	"push2":   OpPush2,
	"push3":   OpPush3,
	"push4":   OpPush4,
	"push5":   OpPush5,
	"push6":   OpPush6,
	"push7":   OpPush7,
	"push8":   OpPush8,
	"push64":  OpPush64,
	"one":     OpOne,
	"true":    OpTrue,
	"neg1":    OpNeg1,
	"pushT":   OpPushT,
	"now":     OpNow,
	"rand":    OpRand,
	"pushL":   OpPushL,
	"add":     OpAdd,
	"sub":     OpSub,
	"mul":     OpMul,
	"div":     OpDiv,
	"mod":     OpMod,
	"not":     OpNot,
	"neg":     OpNeg,
	"inc":     OpInc,
	"dec":     OpDec,
	"index":   OpIndex,
	"len":     OpLen,
	"append":  OpAppend,
	"extend":  OpExtend,
	"slice":   OpSlice,
	"field":   OpField,
	"fieldL":  OpFieldL,
	"ifz":     OpIfz,
	"ifnz":    OpIfnz,
	"else":    OpElse,
	"end":     OpEnd,
	"sum":     OpSum,
	"avg":     OpAvg,
	"max":     OpMax,
	"min":     OpMin,
	"choice":  OpChoice,
	"wChoice": OpWChoice,
	"sort":    OpSort,
	"lookup":  OpLookup,
}

func miniAsm(s string) []Opcode {
	wsp := regexp.MustCompile("[ \t\r\n]")
	words := wsp.Split(strings.ToLower(s), -1)
	opcodes := []Opcode{0}
	for _, w := range words {
		if op, ok := opcodeMap[w]; ok {
			opcodes = append(opcodes, op)
			continue
		}
		// otherwise it should be a hex value
		b, err := strconv.ParseUint(w, 16, 8)
		if err != nil {
			panic(err)
		}
		opcodes = append(opcodes, Opcode(b))
	}
	return opcodes
}

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

func TestMiniAsm(t *testing.T) {
	ops := miniAsm("neg1 zero one push1 45 push2 01 01")
	bytes := []Opcode{0, OpNeg1, OpZero, OpOne, OpPush1, 69, OpPush2, 1, 1}
	assert.Equal(t, ops, bytes)
}

func TestNop(t *testing.T) {
	vm := buildVM(t, "nop")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	assert.Equal(t, vm.Stack().Depth(), 0)
}

func TestPush(t *testing.T) {
	vm := buildVM(t, "neg1 zero one push1 45 push2 01 02 ret")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 1, 69, 513)
}

func TestBigPush(t *testing.T) {
	vm := buildVM(t, "push3 1 2 3 push7 1 2 3 4 5 6 7 push8 fb ff ff ff ff ff ff ff")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 197121, 1976943448883713, -5)
}

func TestPush64(t *testing.T) {
	vm := buildVM(t, "push64 1 2 3 4 5 6 7 8")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	v, err := vm.Stack().Pop()
	assert.Nil(t, err)
	assert.IsType(t, NewID(0), v)
	assert.Equal(t, NewID(578437695752307201), v)
}

func TestDrop(t *testing.T) {
	vm := buildVM(t, "push1 7 nop one zero neg1 drop drop2")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 7)
}

func TestDup(t *testing.T) {
	vm := buildVM(t, "one push1 2 dup push1 3 dup2")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 1, 2, 2, 3, 2, 3)
}

func TestSwapOverPickRoll(t *testing.T) {
	vm := buildVM(t, "zero one push1 2 push1 3 swap over pick 4 roll 4")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 3, 2, 3, 0, 1)
}

func TestMath(t *testing.T) {
	vm := buildVM(t, "push1 55 dup dup add sub push1 7 push1 6 mul dup push1 3 div dup push1 3 mod")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -85, 42, 14, 2)
}

func TestNotNegIncDec(t *testing.T) {
	vm := buildVM(t, "push1 7 not dup not push1 8 neg push1 4 inc push1 6 dec")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 1, -8, 5, 5)
}

func TestIf1(t *testing.T) {
	vm := buildVM(t, "zero ifz push1 13 else push1 42 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf2(t *testing.T) {
	vm := buildVM(t, "zero ifnz push1 13 else push1 42 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf3(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 else push1 42 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 66, 17)
}

func TestIf4(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 else push1 42 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf5(t *testing.T) {
	vm := buildVM(t, "zero ifz push1 13 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIf6(t *testing.T) {
	vm := buildVM(t, "zero ifnz push1 13 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf7(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIf8(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 17)
}

func TestIfNested1(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 zero ifz push1 15 else push1 13 end end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 21, 17)
}

func TestIfNested2(t *testing.T) {
	vm := buildVM(t, "one ifz push1 13 zero ifz push1 15 else push1 13 end end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 17)
}

func TestIfNested3(t *testing.T) {
	vm := buildVM(t, "one ifnz push1 13 zero ifnz push1 15 else push1 13 end end push1 11")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 19, 19, 17)
}
