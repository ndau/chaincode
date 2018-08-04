package vm

import (
	"fmt"
	"strings"
	"time"
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

// Instruction is an opcode with all of its associated data bytes
type Instruction []Opcode

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

// Stats tracks the statistics for a VM
type Stats struct {
	inittime time.Duration // how long it took to call init,
	loadtime time.Duration // how long it took to load and validate the code
	runtime  time.Duration // how long it took to run (not recorded in debug mode)
	maxstack int           // highwater of the stack depth
	childmax int           // highwater of child stacks
}

// AddChildStats is intended to appropriately increment stats based on running child functions
func (vm *ChaincodeVM) AddChildStats(child *ChaincodeVM) {
	vm.stats.inittime += child.stats.inittime
	vm.stats.loadtime += child.stats.loadtime
	// runtime is already accounted for
	// track maximum child depth independently and recursively
	childstacksum := child.stats.maxstack + child.stats.childmax
	if childstacksum > vm.stats.childmax {
		vm.stats.childmax = childstacksum
	}
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
	stats    Stats
}

// New creates a new VM and loads a ChasmBinary into it (or errors)
func New(bin ChasmBinary) (*ChaincodeVM, error) {
	starttime := time.Now()
	vm := ChaincodeVM{}
	defer func() {
		vm.stats.loadtime = time.Now().Sub(starttime)
	}()
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

	// now generate a bitset of used opcodes from the instructions
	usedOpcodes := getUsedOpcodes(generateInstructions(cb.Data[1:]))
	// if it's not a proper subset of the enabled opcodes, don't let it run
	if !usedOpcodes.IsSubsetOf(EnabledOpcodes) {
		return ValidationError{"code contains illegal opcodes"}
	}

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
	starttime := time.Now()
	defer func() {
		vm.stats.inittime = time.Now().Sub(starttime)
	}()
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
	starttime := time.Now()
	defer func() {
		// only record the time for non-debug runs
		if !debug {
			vm.stats.runtime = time.Now().Sub(starttime)
		}
	}()
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
		if d := vm.stack.Depth(); d > vm.stats.maxstack {
			vm.stats.maxstack = d
		}
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

// String implements Stringer so we can print a VM and get something meaningful.
func (vm *ChaincodeVM) String() string {
	st := strings.Split(vm.stack.String(), "\n")
	st1 := make([]string, len(st))
	for i := range st {
		st1[i] = st[i][4:]
	}
	disasm, _ := vm.Disassemble(vm.pc)
	return fmt.Sprintf("%-40s STK: %s\n", disasm, strings.Join(st1, ", "))
}
