package vm

import (
	"fmt"
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

// Chaincode is the type for the VM bytecode program
type Chaincode []Opcode

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

// Randomer is an interface for a type that generates "random" integers (which may vary
// depending on context)
type Randomer interface {
	RandInt() (int64, error)
}

// Nower is an interface for a type that returns the "current" time as a Timestamp
// The definition of "now" may be defined by context.
type Nower interface {
	Now() (Timestamp, error)
}

type funcInfo struct {
	offset int
	nargs  int
}

// ChaincodeVM is the reason we're here
type ChaincodeVM struct {
	runstate  RunState
	code      Chaincode
	stack     *Stack
	pc        int // program counter
	history   []HistoryState
	infunc    int          // the number of the func we're currently in
	handlers  map[byte]int // byte offsets of the handlers by handler ID
	functions []funcInfo   // info for the functions indexed by function number
	rand      Randomer
	now       Nower
}

// New creates a new VM and loads a ChasmBinary into it (or errors)
func New(bin ChasmBinary) (*ChaincodeVM, error) {
	vm := ChaincodeVM{}
	if err := vm.PreLoad(bin); err != nil {
		return nil, err
	}
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

// CreateForFunc creates a new VM from this one that is used to run a function.
// We assume the function number has already been validated.
// and is already in an initialized state to run that function.
// Just call Run() on the new VM after this.
func (vm *ChaincodeVM) CreateForFunc(funcnum int) (*ChaincodeVM, error) {
	finfo := vm.functions[funcnum]
	newstack, err := vm.stack.TopN(finfo.nargs)
	if err != nil {
		return nil, err
	}
	newvm := ChaincodeVM{
		code:      vm.code,
		runstate:  vm.runstate,
		handlers:  vm.handlers,
		functions: vm.functions,
		history:   []HistoryState{},
		infunc:    funcnum,
		pc:        finfo.offset,
		stack:     newstack,
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

// HandlerIDs returns a sorted list of handler IDs that are
// defined for this VM.
func (vm *ChaincodeVM) HandlerIDs() []int {
	ids := []int{}
	for h := range vm.handlers {
		ids = append(ids, int(h))
	}
	sort.Sort(sort.IntSlice(ids))
	return ids
}

// PreLoad is the validation function called before loading a VM to make sure it
// has a hope of loading properly
func (vm *ChaincodeVM) PreLoad(cb ChasmBinary) error {
	return vm.PreLoadOpcodes(cb.Data)
}

// ConvertToOpcodes accepts an array of bytes and returns a Chaincode (array of opcodes)
func ConvertToOpcodes(b []byte) Chaincode {
	ops := make([]Opcode, len(b))
	for i := range b {
		ops[i] = Opcode(b[i])
	}
	return Chaincode(ops)
}

// IsValidChaincode tests if an array of opcodes is a potentially valid
// Chaincode program.
func IsValidChaincode(data Chaincode) error {
	if data == nil {
		return ValidationError{"missing code"}
	}
	if len(data) > maxCodeLength {
		return ValidationError{"code is too long"}
	}
	// make sure the executable part of the code is valid
	_, _, err := validateStructure(data)
	if err != nil {
		return err
	}
	// now generate a bitset of used opcodes from the instructions
	usedOpcodes := getUsedOpcodes(generateInstructions(data))
	// if it's not a proper subset of the enabled opcodes, don't let it run
	if !usedOpcodes.IsSubsetOf(EnabledOpcodes) {
		return ValidationError{"code contains illegal opcodes"}
	}
	return nil
}

// PreLoadOpcodes acepts an array of opcodes and validates it.
// If it fails to validate, the VM is not modified.
// However, if it does validate the VM is updated with
// code and function tables.
func (vm *ChaincodeVM) PreLoadOpcodes(data Chaincode) error {
	if err := IsValidChaincode(data); err != nil {
		return err
	}

	// we know this works because it has already been run
	handlers, functions, _ := validateStructure(data)
	vm.functions = functions
	vm.handlers = handlers
	vm.code = data
	return nil
}

// Init is called to set up the VM to run the handler for a given eventID.
// It can take an arbitrary list of values to push on the stack, which it pushes
// in order -- so if you want something on top of the stack, put it last
// in the argument list. If the VM doesn't have a handler for the specified eventID,
// and it also doesn't have a handler for event 0, then Init will return an error.
func (vm *ChaincodeVM) Init(eventID byte, values ...Value) error {
	stk := NewStack()
	for _, v := range values {
		stk.Push(v)
	}
	return vm.InitFromStack(eventID, stk)
}

// InitFromStack initializes a vm with a given starting stack, which
// should be a new stack
func (vm *ChaincodeVM) InitFromStack(eventID byte, stk *Stack) error {
	vm.stack = stk
	vm.history = []HistoryState{}
	vm.runstate = RsReady
	h, ok := vm.handlers[eventID]
	if !ok {
		h, ok = vm.handlers[0]
		if !ok {
			return ValidationError{"code does not have a handler for the specified event or a default handler"}
		}
	}
	vm.pc = h
	vm.infunc = -1 // we're not in a function to start
	return nil
}

// IP fetches the current instruction pointer (aka program counter)
func (vm *ChaincodeVM) IP() int {
	return vm.pc
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

// Disassemble returns a single disassembled instruction as a text string, possibly with embedded newlines,
// along with how many bytes it consumed.
func (vm *ChaincodeVM) Disassemble(pc int) (string, int) {
	if pc >= len(vm.code) {
		return "END", 0
	}
	op := vm.code[pc]
	numExtra := extraBytes(vm.code, pc)

	out := fmt.Sprintf("%02x:  ", pc)
	sa := []string{fmt.Sprintf("%02x", byte(op))}
	for i := 1; i <= numExtra; i++ {
		sa = append(sa, fmt.Sprintf("%02x", byte(vm.code[pc+i])))
	}
	hex := strings.Join(sa, " ")
	for i := 1; len(hex) > 3*8; i++ {
		out += fmt.Sprintf("%-24s\n%02x:  ", hex[:24], pc+8*i)
		hex = hex[24:]
	}
	args := ""
	if numExtra > 0 && numExtra < 5 {
		args = hex[3:]
	}
	if numExtra >= 5 {
		args = "..."
	}
	out += fmt.Sprintf("%-24s  %-7s %-12s ", hex, op, args)

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
