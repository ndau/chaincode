package vm

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
)

// TODO: tweak data types to support our real keys and timestamps and use ndaumath
//       resolve duration as uint64 or int64
// TODO: calculate and track some execution cost metric
// TODO: test more error states
// TODO: add logging

// The VM package implements a virtual machine for chaincode.

// maxCodeLength is the maximum number of bytes that a VM may contain.
var maxCodeLength = 256

// SetMaxCodeLength allows globally setting the maximum number of bytes a VM may contain.
func SetMaxCodeLength(n int) {
	// TODO: LOG THIS EVENT!
	maxCodeLength = n
}

// Chaincode defines the contract for the virtual machine
type Chaincode interface {
	PreLoad(cb ChasmBinary) error // validates that the code to be loaded is compatible with its stated context
	Init(values []Value)
	Run() (Value, error)
}

// RunState is the current run state of the VM
type RunState byte

// These are runstate constants
const (
	RsNotReady RunState = iota
	RsReady    RunState = iota
	RsRunning  RunState = iota
	RsComplete RunState = iota
	RsError    RunState = iota
)

// HistoryState is a single item in the history of a VM
type HistoryState struct {
	PC    int
	Stack *Stack
	// lists []List
}

// Randomer is an interface for a type that generates "random" integers (which may be
// defined by context)
type Randomer interface {
	RandInt() (int64, error)
}

// Nower is an interface for a type that returns the "current" time as a Timestamp
// The definition of "now" may be defined by context.
type Nower interface {
	Now() (Timestamp, error)
}

// ChaincodeVM is the reason we're here
type ChaincodeVM struct {
	runstate RunState
	context  Opcode
	code     []Opcode
	stack    *Stack
	pc       int // program counter
	history  []HistoryState
	infunc   int   // the number of the func we're currently in
	offsets  []int // byte offsets of the functions
	rand     Randomer
	now      Nower
}

// New creates a new VM and loads a ChasmBinary into it (or errors)
func New(bin ChasmBinary) (*ChaincodeVM, error) {
	vm := ChaincodeVM{}
	if err := vm.PreLoad(bin); err != nil {
		return nil, err
	}
	vm.context = bin.Data[0]
	vm.code = bin.Data[1:]
	vm.runstate = RsNotReady // not ready to run until we've called Init
	r, err := NewDefaultRand()
	if err != nil {
		return nil, err
	}
	vm.rand = r
	n, err := NewDefaultNow()
	if err != nil {
		return nil, err
	}
	vm.now = n
	return &vm, nil
}

// SetRand sets the randomer to call for this VM
func (vm *ChaincodeVM) SetRand(r Randomer) {
	vm.rand = r
}

// SetNow sets the Nower to call for this VM
func (vm *ChaincodeVM) SetNow(n Nower) {
	vm.now = n
}

// CreateForFunc creates a new VM from this one that is used to run a function
func (vm *ChaincodeVM) CreateForFunc(funcnum int, newpc int, nstack int) (*ChaincodeVM, error) {
	newstack, err := vm.stack.TopN(nstack)
	if err != nil {
		return nil, err
	}
	newvm := ChaincodeVM{
		context:  vm.context,
		code:     vm.code,
		runstate: vm.runstate,
		offsets:  vm.offsets,
		history:  []HistoryState{},
		infunc:   funcnum,
		pc:       newpc,
		stack:    newstack,
	}
	return &newvm, nil
}

// extraBytes returns the number of extra bytes associated with a given opcode
func extraBytes(code []Opcode, offset int) int {
	numExtra := 0
	op := code[offset]
	switch op {
	case OpPush1, OpPick, OpRoll, OpDef, OpField, OpFieldL:
		numExtra = 1
	case OpPush2, OpCall, OpDeco:
		numExtra = 2
	case OpPush3:
		numExtra = 3
	case OpPush4:
		numExtra = 4
	case OpPush5:
		numExtra = 5
	case OpPush6:
		numExtra = 6
	case OpPush7:
		numExtra = 7
	case OpPush8, OpPushT:
		numExtra = 8
	case OpPushA, OpPushB:
		numExtra = int(code[offset+1]) + 1
	}
	return numExtra
}

// this validates a VM to make sure that its structural elements (if, else, def, end)
// are well-formed.
// It uses a state machine to generate transitions that depend on the current state and
// the opcode it's evaluating.

type structureState int

// These are the states in the state machine
// The Plus and Minus are temporary states to adjust nesting
const (
	StNull    structureState = iota
	StDef     structureState = iota
	StDefPlus structureState = iota
	StIf      structureState = iota
	StIfPlus  structureState = iota
	StIfMinus structureState = iota
	StElse    structureState = iota
	StError   structureState = iota
)

type tr struct {
	current structureState
	opcode  Opcode
}

// validateStructure reads the script and checks to make sure that the nested
// elements are properly nested and not out of order or missing.
func validateStructure(code []Opcode) ([]int, error) {
	transitions := map[tr]structureState{
		tr{StNull, OpDef}:    StDefPlus,
		tr{StNull, OpIfz}:    StError,
		tr{StNull, OpIfnz}:   StError,
		tr{StNull, OpElse}:   StError,
		tr{StNull, OpEndIf}:  StError,
		tr{StNull, OpEndDef}: StError,
		tr{StDef, OpDef}:     StError,
		tr{StDef, OpIfz}:     StIfPlus,
		tr{StDef, OpIfnz}:    StIfPlus,
		tr{StDef, OpElse}:    StError,
		tr{StDef, OpEndIf}:   StError,
		tr{StDef, OpEndDef}:  StNull,
		tr{StIf, OpDef}:      StError,
		tr{StIf, OpIfz}:      StIfPlus,
		tr{StIf, OpIfnz}:     StIfPlus,
		tr{StIf, OpElse}:     StElse,
		tr{StIf, OpEndIf}:    StIfMinus,
		tr{StIf, OpEndDef}:   StError,
		tr{StElse, OpDef}:    StError,
		tr{StElse, OpIfz}:    StIfPlus,
		tr{StElse, OpIfnz}:   StIfPlus,
		tr{StElse, OpElse}:   StError,
		tr{StElse, OpEndIf}:  StIfMinus,
		tr{StElse, OpEndDef}: StError,
	}

	var state structureState
	var depth int
	var nfuncs int
	var skipcount int
	var offsets = []int{}

	for offset, b := range code {
		// some opcodes have operands and we need to be sure those are skipped
		if skipcount > 0 {
			skipcount--
			continue
		}
		skipcount = extraBytes(code, offset)
		newstate, found := transitions[tr{state, b}]
		if !found {
			continue
		}
		switch newstate {
		case StDefPlus:
			funcnum := int(code[offset+1])
			if funcnum != nfuncs {
				return offsets, ValidationError{fmt.Sprintf("def should have been %d, found %d", nfuncs, funcnum)}
			}
			offsets = append(offsets, int(offset+2)) // skip the def opcode, which is 2 bytes
			nfuncs++
			newstate = StDef
		case StError:
			return offsets, ValidationError{fmt.Sprintf("invalid structure at offset %d", offset)}
		case StIfPlus:
			depth++
			newstate = StIf
		case StIfMinus:
			depth--
			if depth == 0 {
				newstate = StDef
			} else {
				newstate = StIf
			}
		}
		state = newstate
	}
	if skipcount != 0 {
		return offsets, ValidationError{"missing operands"}
	}
	if state != StNull {
		return offsets, ValidationError{"missing end"}
	}
	if nfuncs < 1 {
		return offsets, ValidationError{"missing def"}
	}
	return offsets, nil
}

// Stack returns the current stack of the VM
func (vm *ChaincodeVM) Stack() *Stack {
	return vm.stack
}

// History returns the entire history of this VM's operation
func (vm *ChaincodeVM) History() []HistoryState {
	return vm.history
}

// PreLoad is the validation function called before loading a VM to make sure it
// has a hope of loading properly
func (vm *ChaincodeVM) PreLoad(cb ChasmBinary) error {
	if cb.Data == nil {
		return ValidationError{"missing code"}
	}
	if len(cb.Data) > maxCodeLength {
		return ValidationError{"code is too long"}
	}
	// make sure the executable part of the code is valid
	offsets, err := validateStructure(cb.Data[1:])
	if err != nil {
		return err
	}
	vm.offsets = offsets

	ctx, ok := ContextLookup(cb.Context)
	if !ok {
		return ValidationError{"invalid context string"}
	}
	if _, ok := Contexts[ContextByte(cb.Data[0])]; !ok {
		return ValidationError{"invalid context byte"}
	}
	if ctx != ContextByte(cb.Data[0]) {
		return ValidationError{"context byte and context string disagree"}
	}
	// we seem to be OK
	return nil
}

// Init is called to set up the VM to run
func (vm *ChaincodeVM) Init(values ...Value) {
	vm.stack = newStack()
	for _, v := range values {
		vm.stack.Push(v)
	}
	vm.history = []HistoryState{}
	vm.runstate = RsReady
	vm.pc = 2 // skip the def 0 at the start
}

// Run runs a VM from its current state until it ends
func (vm *ChaincodeVM) Run(debug bool) error {
	if debug {
		vm.DisassembleAll()
	}
	if vm.runstate == RsReady {
		vm.runstate = RsRunning
	}
	for vm.runstate == RsRunning {
		if debug {
			fmt.Println(vm)
		}
		if err := vm.Step(debug); err != nil {
			return err
		}
	}
	return nil
}

func (vm *ChaincodeVM) runtimeError(err error) error {
	rte := wrapRuntimeError(err)
	return rte.PC(vm.pc - 1)
}

// This is only run on VMs that have been validated
func (vm *ChaincodeVM) skipToMatchingBracket() error {
	for {
		instr := vm.code[vm.pc]
		vm.pc++
		nesting := 0
		switch instr {
		case OpIfnz, OpIfz:
			nesting++
		case OpElse:
			if nesting == 0 {
				// we're at the right level, so we're done
				return nil
			}
		case OpEndIf:
			if nesting > 0 {
				nesting--
			} else {
				// we're at the right level so we're done
				return nil
			}
		default:
			// fail-safe (should never happen)
			if vm.pc > len(vm.code) {
				return vm.runtimeError(newRuntimeError("VM RAN OFF THE END!"))
			}
		}
	}
}

// callFunction calls the function numbered by funcnum, copying nargs to a new stack
// returns the value left on the stack by the called function
func (vm *ChaincodeVM) callFunction(funcnum int, nargs int, debug bool, extraArgs ...Value) (Value, error) {
	var retval Value
	if funcnum <= vm.infunc || funcnum >= len(vm.offsets) {
		return retval, vm.runtimeError(newRuntimeError("invalid function number (no recursion allowed)"))
	}
	newpc := vm.offsets[funcnum]

	childvm, err := vm.CreateForFunc(funcnum, newpc, nargs)
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	for _, e := range extraArgs {
		err := childvm.stack.Push(e)
		if err != nil {
			return retval, vm.runtimeError(err)
		}
	}
	err = childvm.Run(debug)
	// no matter what, we want the history
	vm.history = append(vm.history, childvm.history...)
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	// we've called the child function, now get back its return value
	retval, err = childvm.stack.Pop()
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	return retval, nil
}

// Step executes a single instruction
func (vm *ChaincodeVM) Step(debug bool) error {
	switch vm.runstate {
	default:
		return newRuntimeError("vm is not in runnable state")
	case RsReady:
		vm.runstate = RsRunning
		fallthrough
	case RsRunning:
		vm.history = append(vm.history, HistoryState{
			PC:    vm.pc,
			Stack: vm.stack.Clone(),
			// lists: vm.lists[:],
		})
	}

	// Check to see if we're still in runnable code
	if vm.pc >= len(vm.code) {
		vm.runstate = RsComplete
		return nil
	}

	// fetch the instruction
	instr := vm.code[vm.pc]
	// we always increment the pc immediately; we may also add to it if an instruction has additional data
	// when we report errors, we subtract 1 to get the correct value
	vm.pc++
	switch instr {
	case OpNop:
		// do nothing
	case OpDrop:
		// discards top element on stack
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
	case OpDrop2:
		// discards top two elements on stack
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
	case OpDup:
		v, err := vm.stack.Peek()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpDup2:
		v1, err := vm.stack.Get(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		v0, err := vm.stack.Get(0)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v0); err != nil {
			return vm.runtimeError(err)
		}
	case OpSwap:
		v1, err := vm.stack.PopAt(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
	case OpOver:
		v1, err := vm.stack.Get(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
	case OpPick:
		n := int(vm.code[vm.pc])
		vm.pc++
		v, err := vm.stack.Get(n)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpRoll:
		n := int(vm.code[vm.pc])
		vm.pc++
		v, err := vm.stack.PopAt(n)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpRet:
		vm.runstate = RsComplete
	case OpFail:
		vm.runstate = RsError
		return vm.runtimeError(newRuntimeError("fail opcode invoked"))
	case OpZero:
		if err := vm.stack.Push(NewNumber(0)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPush1, OpPush2, OpPush3, OpPush4, OpPush5, OpPush6, OpPush7, OpPush8:
		// use a mask to retrieve the actual count of bytes to fetch
		nbytes := byte(instr) & 0xF
		var value int64
		var i byte
		var b Opcode
		for i = 0; i < nbytes; i++ {
			b = vm.code[vm.pc]
			vm.pc++
			value |= int64(b) << (i * 8)
		}
		if b&0x80 != 0 {
			for i := nbytes; i < 8; i++ {
				value |= 0xFF
			}
		}
		if err := vm.stack.Push(NewNumber(value)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushT:
		var value int64
		var i byte
		var b Opcode
		for i = 0; i < 8; i++ {
			b = vm.code[vm.pc]
			vm.pc++
			value |= int64(b) << (i * 8)
		}
		ts := NewTimestamp(value)
		if err := vm.stack.Push(ts); err != nil {
			return vm.runtimeError(err)
		}
	case OpNow:
		ts, err := vm.now.Now()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(ts); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushB, OpPushA:
		n := int(vm.code[vm.pc])
		vm.pc++
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte(vm.code[vm.pc])
			vm.pc++
		}
		v := NewBytes(b)
		if err := vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpOne:
		if err := vm.stack.Push(NewNumber(1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNeg1:
		if err := vm.stack.Push(NewNumber(-1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpRand:
		r, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(r)); err != nil {
			return vm.runtimeError(err)
		}

	case OpPushL:
		if err := vm.stack.Push(NewList()); err != nil {
			return vm.runtimeError(err)
		}
	case OpAdd, OpMul, OpDiv, OpMod:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch instr {
		case OpAdd:
			t = n2 + n1
		case OpMul:
			t = n2 * n1
		case OpDiv:
			if n1 == 0 {
				return vm.runtimeError(newRuntimeError("divide by zero"))
			}
			t = n2 / n1
		case OpMod:
			if n1 == 0 {
				return vm.runtimeError(newRuntimeError("divide by zero"))
			}
			t = n2 % n1
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpSub:
		// Subtraction is special because you can also subtract timestamps
		v1, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		v2, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch n1 := v1.(type) {
		case Number:
			n2, ok := v2.(Number)
			if !ok {
				return vm.runtimeError(newRuntimeError("incompatible types"))
			}
			t = n2.v - n1.v
		case Timestamp:
			n2, ok := v2.(Timestamp)
			if !ok {
				return vm.runtimeError(newRuntimeError("incompatible types"))
			}
			t = int64(n2.t) - int64(n1.t)
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNot, OpNeg, OpInc, OpDec:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch instr {
		case OpNot:
			if n1 == 0 {
				t = 1
			}
		case OpNeg:
			t = -n1
		case OpInc:
			t = n1 + 1
		case OpDec:
			t = n1 - 1
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpEq, OpGt, OpLt:
		v2, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		v1, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		r, err := v1.Compare(v2)
		if err != nil {
			return vm.runtimeError(err)
		}
		var result int64
		switch instr {
		case OpEq:
			if r == 0 {
				result = 1
			}
		case OpGt:
			if r > 0 {
				result = 1
			}
		case OpLt:
			if r < 0 {
				result = 1
			}
		}
		if err := vm.stack.Push(NewNumber(result)); err != nil {
			return vm.runtimeError(err)
		}

	case OpIndex:
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if n >= l.Len() {
			return vm.runtimeError(newRuntimeError("list index out of bounds"))
		}
		if err := vm.stack.Push(l[n]); err != nil {
			return vm.runtimeError(err)
		}

	case OpLen:
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(int64(len(l)))); err != nil {
			return vm.runtimeError(err)
		}

	case OpAppend:
		v, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		// Limit list size
		if l.Len()+1 > MaxListSize {
			return vm.runtimeError(newRuntimeError("resulting list too large"))
		}
		newlist := l.Append(v)
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpExtend:
		l1, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		l2, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		// Limit list size
		if l1.Len()+l2.Len() > MaxListSize {
			return vm.runtimeError(newRuntimeError("resulting list too large"))
		}

		newlist := l2.Extend(l1)
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpSlice:
		end, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		begin, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if begin < 0 || begin > l.Len() || end < 0 || end > l.Len() || begin > end {
			return vm.runtimeError(newRuntimeError("index out of range in slice"))
		}
		newlist := l[begin:end]
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpField:
		st, err := vm.stack.PopAsStruct()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		f, err := st.Field(int(fix))
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(f); err != nil {
			return vm.runtimeError(err)
		}

	case OpFieldL:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		extract := func(v Value) (Value, error) {
			return v.(Struct).Field(int(fix))
		}
		result, err := src.Map(extract)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpDef:
		// if we try to execute a def statement there has been an error and we should abort
		return vm.runtimeError(newRuntimeError("tried to execute def opcode"))

	case OpCall:
		// The call opcode tracks the number of the current routine, and will only call a
		// function that has a larger number than itself (this prevents recursion). It constructs a
		// new stack for the function and populates it by copying (NOT popping off!) the specified number of
		// values from the caller's stack. The function call returns a single Value which is pushed
		// onto the caller's stack.
		funcnum := int(vm.code[vm.pc])
		nargs := int(vm.code[vm.pc+1])
		vm.pc += 2
		result, err := vm.callFunction(funcnum, nargs, debug)
		if err != nil {
			return err
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpDeco:
		funcnum := int(vm.code[vm.pc])
		nargs := int(vm.code[vm.pc+1])
		vm.pc += 2
		// we're going to iterate over a List
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		newlist := NewList()
		for i := range l {
			s, ok := l[i].(Struct)
			if !ok {
				return vm.runtimeError(newRuntimeError("list element should have been struct"))
			}

			retval, err := vm.callFunction(funcnum, nargs, debug, s)
			if err != nil {
				return err
			}
			// in order to prevent memory bombs, deco cannot add non-scalars
			if !retval.IsScalar() {
				return vm.runtimeError(newRuntimeError("deco result must be scalar"))
			}
			newlist = newlist.Append(s.Append(retval))
		}
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpEndDef:
		// we hit this at the end of a function that hasn't used OpRet or OpFail
		vm.runstate = RsComplete

	case OpIfz:
		t, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		success := false
		switch v := t.(type) {
		case Number:
			success = (v.v == 0)
		default:
			// nothing else can test true for zeroness
		}
		// if we did not succeed, we have to skip to the matching end or else opcode
		if !success {
			if err := vm.skipToMatchingBracket(); err != nil {
				return vm.runtimeError(err)
			}
		}
	case OpIfnz:
		t, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		success := false
		switch v := t.(type) {
		case Number:
			success = (v.v != 0)
		default:
			// everything else tests true for nonzeroness
			success = true
		}
		// if we did not succeed, we have to skip to the matching END or ELSE opcode
		if !success {
			if err := vm.skipToMatchingBracket(); err != nil {
				return vm.runtimeError(err)
			}
		}
	case OpElse:
		// if we hit this in execution, it means we did the first clause of an if statement
		// and now need to skip to the matching end
		if err := vm.skipToMatchingBracket(); err != nil {
			return vm.runtimeError(err)
		}
	case OpEndIf:
		// OpEndIf is a no-op (it is only hit when it ends an if block that evaluated to true
		// and there was no Else clause

	case OpSum, OpAvg, OpMax, OpMin:
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}

		// Define helper functions for the Reduce function
		sum := func(prev, current Value) Value {
			p := prev.(Number).v
			c := current.(Number).v
			return NewNumber(p + c)
		}
		max := func(prev, current Value) Value {
			cmp, _ := prev.Compare(current)
			if cmp < 0 {
				return current
			}
			return prev
		}
		min := func(prev, current Value) Value {
			cmp, _ := current.Compare(prev)
			if cmp < 0 {
				return current
			}
			return prev
		}

		var result Value
		switch instr {
		case OpSum:
			result = l.Reduce(sum, NewNumber(0))
		case OpAvg:
			if l.Len() == 0 {
				return vm.runtimeError(newRuntimeError("cannot average empty list"))
			}
			result = NewNumber(l.Reduce(sum, NewNumber(0)).(Number).v / l.Len())
		case OpMin:
			result = l.Reduce(min, NewNumber(math.MaxInt64))
		case OpMax:
			result = l.Reduce(max, NewNumber(math.MinInt64))
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpChoice:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if src.Len() == 0 {
			return vm.runtimeError(newRuntimeError("cannot use choice on empty list"))
		}
		i, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}
		r := rand.New(rand.NewSource(i))
		n := r.Intn(int(src.Len()))
		item := src[n]
		if err := vm.stack.Push(item); err != nil {
			return vm.runtimeError(err)
		}

	case OpWChoice:
		fix := vm.code[vm.pc]
		vm.pc++
		src, err := vm.stack.PopAsListOfStructs(int(fix))
		if err != nil {
			return vm.runtimeError(err)
		}
		// because of PopAsListOfStructs(), we know we're safe to traverse
		// the list of structs and pull out our specified field as a Number
		sum := func(prev, current Value) Value {
			p := prev.(Number).v
			fi, _ := current.(Struct).Field(int(fix))
			c := fi.(Number).v
			return NewNumber(p + c)
		}
		total := src.Reduce(sum, NewNumber(0)).(Number).AsInt64()

		rand, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}

		var partialSum int64
		for i := range src {
			fi, _ := src[i].(Struct).Field(int(fix))
			partialSum += fi.(Number).AsInt64()
			if FractionLess(rand, math.MaxInt64, partialSum, total) {
				err := vm.stack.Push(src[i])
				if err != nil {
					return vm.runtimeError(err)
				}
				return nil
			}
		}

		// if we get here, something is very wrong
		return vm.runtimeError(newRuntimeError(fmt.Sprintf("wchoice can't happen: %d %d %d", rand, partialSum, total)))

	case OpSort:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		// note - error handling is weak because the less function cannot fail
		// so we can only figure it out after the sort completes
		hadErr := false
		less := func(i, j int) bool {
			fi, e1 := src[i].(Struct).Field(int(fix))
			fj, e2 := src[j].(Struct).Field(int(fix))
			cmp, e3 := fi.Compare(fj)
			if e1 != nil || e2 != nil || e3 != nil {
				hadErr = true
			}
			return cmp == -1
		}
		sort.Slice(src, less)
		if hadErr {
			return vm.runtimeError(newRuntimeError("sort error"))
		}
		if err := vm.stack.Push(src); err != nil {
			return vm.runtimeError(err)
		}

	// case OpLookup:
	default:
		return vm.runtimeError(newRuntimeError("unimplemented opcode"))
	}

	return nil
}

// DisassembleAll dumps a disassembly of the whole VM
func (vm *ChaincodeVM) DisassembleAll() {
	fmt.Println("--DISASSEMBLY--")
	for pc := 0; pc < len(vm.code); {
		s, delta := vm.Disassemble(pc)
		pc += delta
		fmt.Println(s)
	}
	fmt.Println("---------------")
}

// Disassemble returns a single disassembled instruction, along with how many bytes it consumed
func (vm *ChaincodeVM) Disassemble(pc int) (string, int) {
	if pc >= len(vm.code) {
		return "END", 0
	}
	op := vm.code[pc]
	numExtra := extraBytes(vm.code, pc)
	sa := []string{fmt.Sprintf("%3d  %02x", pc, byte(op))}
	for i := numExtra; i > 0; i-- {
		sa = append(sa, fmt.Sprintf("%02x", byte(vm.code[pc+i])))
	}
	hex := strings.Join(sa, " ")
	if len(hex) > 30 {
		hex = hex[:25] + "..."
	}
	out := fmt.Sprintf("%-30s  %s", hex, op)

	return out, numExtra + 1
}

func (vm *ChaincodeVM) String() string {
	st := strings.Split(vm.stack.String(), "\n")
	st1 := make([]string, len(st))
	for i := range st {
		st1[i] = st[i][4:]
	}
	disasm, _ := vm.Disassemble(vm.pc)
	return fmt.Sprintf("%-40s STK: %s\n", disasm, strings.Join(st1, ", "))
}
