package vm

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeco1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 deco 0 0 fieldl 2 sum enddef
		def 0 dup field 0 dup mul swap  field 1 dup mul add enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 455)
}

func TestStringers(t *testing.T) {
	assert.Equal(t, "Call", OpCall.String())
	vid := NewBytes([]byte("hi"))
	assert.Equal(t, "hi", vid.String())
	vn := NewNumber(123)
	assert.Equal(t, "123", vn.String())
	vt := NewTimestamp(0)
	assert.Equal(t, "2000-01-01T00:00:00Z", vt.String())
	vl := NewList()
	vl = vl.Append(NewBytes([]byte("July"))).Append(NewNumber(18))
	assert.Equal(t, "[July, 18]", vl.String())
	vs := NewStruct(NewBytes([]byte("July")), NewNumber(18))
	assert.Equal(t, "str(0)[July, 18]", vs.String())
}

func TestExerciseStrings(t *testing.T) {
	vm := buildVM(t, "handler 0 sort 6 push1 3 index field 1 enddef")
	vm.Init(0)

	assert.Contains(t, vm.String(), "Sort")
	da, n := vm.Disassemble(4)
	assert.Equal(t, 2, n)
	assert.Contains(t, da, "Push1")
}

func TestLookup1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 0 push1 4 gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestLookup2(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 1 push1 4 gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 2)
}

func TestLookupFail1(t *testing.T) {
	vm := buildVM(t, `
		handler 0 lookup 0 0 enddef
		def 0 field 1 push1 FF gt enddef
	`)
	l := NewList()
	for i := int64(0); i < 5; i++ {
		st := NewStruct(NewNumber(2*i), NewNumber(3*i+1))
		l = l.Append(st)
	}
	vm.Init(0, l)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestUnimplemented(t *testing.T) {
	// first make sure that the validation check forbids an invalid opcode
	buildVMfail(t, "handler 0 FF enddef")

	// now let's hack a VM after it passes validation to contain illegal data
	vm := buildVM(t, "handler 0 NOP enddef")
	// replace the nop with FF and try to run it; should still fail
	vm.code[3] = Opcode(0xFF)
	vm.Init(0)
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestUnderflows(t *testing.T) {
	p := regexp.MustCompile("[[:space:]]+")
	keywords := p.Split(`drop drop2 dup dup2 swap over
		add sub mul div mod divmod muldiv not neg inc dec
		eq lt gt index len append extend slice sum avg max min`, -1)
	for _, k := range keywords {
		prog := "handler 0 " + k + " enddef"
		vm := buildVM(t, prog)
		vm.Init(0)
		err := vm.Run(false)
		assert.NotNil(t, err)
		correct := strings.HasPrefix(err.Error(), "stack underflow") ||
			strings.HasPrefix(err.Error(), "stack index error")
		assert.True(t, correct, "Keyword=%s msg=%s", k, err)
	}
}

func TestDisableOpcode(t *testing.T) {
	// now let's hack a VM after it passes validation to contain illegal data
	vm := buildVM(t, "handler 0 NOP enddef")
	vm.Init(0)
	err := vm.Run(false)
	assert.Nil(t, err)

	DisableOpcode(OpNop)
	// now the validation check should fail an invalid opcode
	buildVMfail(t, "handler 0 NOP enddef")
	// but we have to re-enable Nop or other tests might fail
	EnabledOpcodes.Set(int(OpNop))
}

func TestNegativeIndex(t *testing.T) {
	prog := `Handler 00
		Neg1 Index
		EndDef`
	vm := buildVM(t, prog)
	vm.Init(0, NewList().Append(NewNumber(1)).Append(NewNumber(2)))
	err := vm.Run(false)
	assert.NotNil(t, err)
}

func TestIndex2(t *testing.T) {
	// this test is making sure that the 8f embedded into the PushB doesn't
	// cause skipToMatchingBracket to fail
	prog := `Handler 00
		IfZ
		PushB 8 b6 42 59 a3 8f 28 81 70
		EndIf
		EndDef`
	vm := buildVM(t, prog)
	vm.Init(0, NewList().Append(NewNumber(1)).Append(NewNumber(2)))
	err := vm.Run(false)
	assert.Nil(t, err)
}

func TestNoHandlers(t *testing.T) {
	// this test is making sure that if no handlers are defined, the vm won't load
	buildVMfail(t, "Def 00 Nop EndDef")
}

func TestMultipleHandlers(t *testing.T) {
	// this tests that we can define and call multiple handlers
	prog := `handler 1 10 add enddef
		handler 1 8 sub enddef
		handler 1 30 mul enddef
		handler 1 3 div enddef
		`
	vm := buildVM(t, prog)

	vm.Init(16, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(8, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 8)

	vm.Init(48, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(3, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 3)
}

func TestDefaultHandler(t *testing.T) {
	// this tests that we can define a different default handler
	// and also that it gets called if we invoke an event not in our list
	prog := `handler 1 10 add enddef
		handler 1 8 sub enddef
		handler 0 mod enddef
		handler 1 30 mul enddef
		handler 1 3 div enddef
		`
	vm := buildVM(t, prog)

	vm.Init(0, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)

	vm.Init(8, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 8)

	vm.Init(77, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 0)
}

func TestMultipleEvents(t *testing.T) {
	// this tests that we can define a different default handler
	// and also that it gets called if we invoke an event not in our list
	prog := `handler 3 10 12 14 add enddef
		handler 2 0 5 mul enddef
		`
	vm := buildVM(t, prog)

	vm.Init(18, NewNumber(12), NewNumber(4))
	err := vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(18, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(20, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 16)

	vm.Init(0, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(5, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)

	vm.Init(77, NewNumber(12), NewNumber(4))
	err = vm.Run(false)
	assert.Nil(t, err)
	checkStack(t, vm.Stack(), 48)
}

func TestHandlerIDs(t *testing.T) {
	// this tests that the HandlerIDs function works right
	prog := `handler 3 10 12 14 add enddef
		handler 2 0 5 mul enddef
		`
	vm := buildVM(t, prog)
	assert.Equal(t, []int{0, 5, 16, 18, 20}, vm.HandlerIDs())

	prog = `handler 1 1 add enddef`
	vm = buildVM(t, prog)
	assert.Equal(t, []int{1}, vm.HandlerIDs())

	prog = `handler 0 add enddef`
	vm = buildVM(t, prog)
	assert.Equal(t, []int{0}, vm.HandlerIDs())
}

func TestNumericBinops(t *testing.T) {
	// This exercises the basic binary operators using a table-driven approach.
	f := func(op Opcode, a, b int64) (int64, error) {
		prog := "handler 0 " + op.String() + " enddef"
		vm := buildVM(t, prog)
		vm.Init(0, NewNumber(a), NewNumber(b))
		err := vm.Run(false)
		if err != nil {
			return 0, err
		}
		return vm.Stack().PopAsInt64()
	}

	tests := []struct {
		name    string
		op      Opcode
		a       int64
		b       int64
		want    int64
		wantErr bool
	}{
		{"add a", OpAdd, 12, 3, 15, false},
		{"sub a", OpSub, 12, 3, 9, false},
		{"sub b", OpSub, 3, 12, -9, false},
		{"mul a", OpMul, 12, 3, 36, false},
		{"mul b", OpMul, 12, -3, -36, false},
		{"mul c", OpMul, -12, 3, -36, false},
		{"mul d", OpMul, -12, -3, 36, false},
		{"div a", OpDiv, 12, 3, 4, false},
		{"div b", OpDiv, 12, -3, -4, false},
		{"div c", OpDiv, -12, 3, -4, false},
		{"div d", OpDiv, -12, -3, 4, false},
		{"div e", OpDiv, 12, 5, 2, false},
		{"div f", OpDiv, 0, 5, 0, false},
		{"div g", OpDiv, 5, 0, 0, true},
		{"mod a", OpMod, 6, 5, 1, false},
		{"mod b", OpMod, 6, -5, 1, false},
		{"mod c", OpMod, -6, 5, -1, false},
		{"mod d", OpMod, -6, -5, -1, false},
		{"mod e", OpMod, 12, 4, 0, false},
		{"mod f", OpMod, 0, 5, 0, false},
		{"mod g", OpMod, 5, 0, 0, true},
		{"mod h", OpMod, 5, 7, 5, false},
		{"or a", OpOr, 0x55, 0x0f, 0x5f, false},
		{"or b", OpOr, 0, -1, -1, false},
		{"or c", OpOr, 0x5555555555555555, ^0x5555555555555555, -1, false},
		{"and a", OpAnd, 0x55, 0x0f, 0x05, false},
		{"and b", OpAnd, 0, -1, 0, false},
		{"and c", OpAnd, 0x5555555555555555, ^0x5555555555555555, 0, false},
		{"xor a", OpXor, 0x55, 0x0f, 0x5a, false},
		{"xor b", OpXor, 0, -1, -1, false},
		{"xor c", OpXor, 0x5555555555555555, ^0x5555555555555555, -1, false},
		{"xor d", OpXor, 0x5555555555555555, -1, ^0x5555555555555555, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f(tt.op, tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Binop error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%v %s %v = %v, want %v", tt.a, tt.op, tt.b, got, tt.want)
			}
		})
	}
}

func TestNumericUnaryOps(t *testing.T) {
	// This exercises the basic unary operators using a table-driven approach.
	f := func(op Opcode, a int64) (int64, error) {
		prog := "handler 0 " + op.String() + " enddef"
		vm := buildVM(t, prog)
		vm.Init(0, NewNumber(a))
		err := vm.Run(false)
		if err != nil {
			return 0, err
		}
		return vm.Stack().PopAsInt64()
	}

	tests := []struct {
		name    string
		op      Opcode
		a       int64
		want    int64
		wantErr bool
	}{
		{"not a", OpNot, 3, 0, false},
		{"not b", OpNot, -1, 0, false},
		{"not c", OpNot, 2313, 0, false},
		{"not d", OpNot, 0, -1, false},
		{"not e", OpNot, 1, 0, false},
		{"bnot a", OpBNot, 0, -1, false},
		{"bnot b", OpBNot, -1, 0, false},
		{"bnot c", OpBNot, 1, -2, false},
		{"bnot d", OpBNot, -3, 2, false},
		{"count1s a", OpCount1s, 3, 2, false},
		{"count1s b", OpCount1s, 0x5A, 4, false},
		{"count1s c", OpCount1s, -1, 64, false},
		{"count1s d", OpCount1s, 0, 0, false},
		{"count1s e", OpCount1s, 0x5555555555555555, 32, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f(tt.op, tt.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%s %v = %v, want %v", tt.op, tt.a, got, tt.want)
			}
		})
	}
}
