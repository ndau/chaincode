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
		assert.Equal(t, n, values[len(values)-i-1])
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
	vm := buildVM(t, "neg1 zero one push1 45 push2 01 01 ret")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), -1, 0, 1, 69, 257)
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

func TestSwapOver(t *testing.T) {
	vm := buildVM(t, "zero one push1 2 push1 3 swap over pick 4 roll 4")
	vm.Init(nil)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0, 3, 2, 3, 0, 1)
}
