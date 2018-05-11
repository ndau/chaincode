package vm

import "fmt"

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
	PreLoad(context byte) error // validates that the code to be loaded is compatible with its stated context
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

type state struct {
	pc    int
	stack *Stack
	// lists []List
}

// ChaincodeVM is the reason we're here
type ChaincodeVM struct {
	runstate RunState
	context  byte
	code     []byte
	stack    *Stack
	// lists    []List
	pc      int // program counter
	history []state
}

// NewVM creates a new VM and loads code into it (or errors)
func NewVM(code []byte) (*ChaincodeVM, error) {
	vm := ChaincodeVM{}
	if err := vm.PreLoad(code); err != nil {
		return nil, err
	}
	vm.context = code[0]
	vm.code = code[1:]
	vm.runstate = RsNotReady // not ready to run until we've called Init
	return &vm, nil
}

// ValidationError is returned when the code is invalid and cannot be loaded or run
type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}

// RuntimeError is returned when the vm encounters an error during execution
type RuntimeError struct {
	pc  int
	msg string
}

// PC sets the program counter value for an error
func (e RuntimeError) PC(pc int) RuntimeError {
	e.pc = pc
	return e
}

func newRuntimeError(s string) error {
	return RuntimeError{pc: -1, msg: s}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[pc=%d] %s", e.pc, e.msg)
}

func validateNesting(code []byte) error {
	nesting := 0
	haselse := []bool{}
	for _, b := range code {
		switch b {
		case OpIfnz, OpIfz:
			nesting++
			haselse = append(haselse, false)
		case OpElse:
			if nesting == 0 {
				return ValidationError{"invalid nesting (else without if)"}
			}
			if haselse[nesting] {
				return ValidationError{"invalid nesting (too many elses)"}
			}
			haselse[nesting] = true
		case OpEnd:
			if nesting == 0 {
				return ValidationError{"invalid nesting (end without if)"}
			}
			nesting--
			haselse = haselse[:len(haselse)-1]
		default:
		}
	}
	if nesting != 0 {
		return ValidationError{"invalid nesting (if without end)"}
	}
	// we think we're ok to try to load this
	return nil
}

// PreLoad is the validation function called before loading a VM to make sure it
// has a hope of loading properly
func (vm *ChaincodeVM) PreLoad(code []byte) error {
	if len(code) == 0 || code == nil {
		return ValidationError{"code is empty"}
	}
	if len(code) > maxCodeLength {
		return ValidationError{"code is too long"}
	}
	switch code[0] {
	case CtxTest:
		return nil
	case CtxEaiTiming:
		return nil
	case CtxMarketPrice:
		return nil
	case CtxNodePayout:
		return nil
	case CtxNodeQuality:
		return nil
	default:
		return ValidationError{"invalid context"}
	}
}

// Init is called to set up the VM to run
func (vm *ChaincodeVM) Init(values []Value) {
	vm.stack = newStack()
	for _, v := range values {
		vm.stack.Push(v)
	}
	// TODO: lists
	// vm.lists = make([]List, 0)
	vm.history = []state{}
	vm.runstate = RsReady
	vm.pc = 0
}

// Run runs a VM from its current state until it ends
func (vm *ChaincodeVM) Run() error {
	if vm.runstate == RsReady {
		vm.runstate = RsRunning
	}
	for vm.runstate == RsRunning {
		if err := vm.Step(); err != nil {
			return err
		}
	}
	return nil
}

func (vm *ChaincodeVM) runtimeError(err error) error {
	if e, ok := err.(RuntimeError); ok {
		return e.PC(vm.pc - 1)
	}
	return err
}

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
				// we're at the right level, so just skip this instruction and continue
				vm.pc++ // skip
				return nil
			}
		case OpEnd:
			if nesting > 0 {
				nesting--
			} else {
				// we're at the right level so we're done
				vm.pc++ // consume this one too
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

// Step executes a single instruction
func (vm *ChaincodeVM) Step() error {
	switch vm.runstate {
	default:
		return newRuntimeError("vm is not in runnable state!")
	case RsReady:
		vm.runstate = RsRunning
		fallthrough
	case RsRunning:
		vm.history = append(vm.history, state{
			pc:    vm.pc,
			stack: vm.stack.Clone(),
			// lists: vm.lists[:],
		})
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
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		v, err := vm.stack.Get(int(n))
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpRoll:
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		v, err := vm.stack.PopAt(int(n))
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
		if err := vm.stack.Push(newNumber(0)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushN + 1, OpPushN + 2, OpPushN + 3, OpPushN + 4, OpPushN + 5, OpPushN + 6, OpPushN + 7:
		// use a mask to retrieve the actual count of bytes to fetch
		nbytes := instr & 0x7
		var value int64
		var b, i byte
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
		if err := vm.stack.Push(newNumber(value)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPush64:
	case OpOne:
		if err := vm.stack.Push(newNumber(1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNeg1:
		if err := vm.stack.Push(newNumber(-1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushT:
	case OpNow:
	case OpRand:
	case OpPushL:
	case OpAdd:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n2 + n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpSub:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n2 - n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpMul:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n2 * n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpDiv:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n2 / n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpMod:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n2 % n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNot:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		if n1 == 0 {
			t = 1
		}
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNeg:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := -n1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpInc:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n1 + 1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpDec:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		t := n1 - 1
		if err := vm.stack.Push(newNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpIndex:
	case OpLen:
	case OpAppend:
	case OpExtend:
	case OpSlice:
	case OpField:
	case OpFieldL:
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
	case OpEnd:
		// this is a nop
	case OpSum:
	case OpAvg:
	case OpMax:
	case OpMin:
	case OpChoice:
	case OpWChoice:
	case OpSort:
	case OpLookup:
	default:
	}

	return nil
}
