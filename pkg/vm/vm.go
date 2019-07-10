package vm

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"

	math "github.com/oneiro-ndev/ndaumath/pkg/types"
)

// TODO: tweak data types to support our real keys and timestamps and use ndaumath
//       resolve duration as uint64 or int64
// TODO: calculate and track some execution cost metric
// TODO: test more error states
// TODO: add logging

// The VM package implements a virtual machine for chaincode.

// maxCodeLength is the maximum number of bytes that a VM may contain, excluding
// the data bytes for PushB.
var maxCodeLength = 256

// maxTotalLength is the maximum number of bytes that a VM may contain in total
var maxTotalLength = 1024

// SetMaxLengths allows globally setting the maximum number of bytes a VM may contain.
func SetMaxLengths(code, total int) {
	maxCodeLength = code
	maxTotalLength = total
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

// DisassembleTabWidth is the number of spaces of indent for each tab in disassembly
const DisassembleTabWidth = 4

var tab string

func init() {
	for i := 0; i < DisassembleTabWidth; i++ {
		tab += " "
	}
}

// IndentOpcodes is the set of opcodes which result in a formatting indent
var IndentOpcodes map[Opcode]struct{}

// DedentOpcodes is the set of opcodes which result in a formatting dedent
var DedentOpcodes map[Opcode]struct{}

func init() {
	s := struct{}{}

	IndentOpcodes = make(map[Opcode]struct{})
	IndentOpcodes[OpDef] = s
	IndentOpcodes[OpIfZ] = s
	IndentOpcodes[OpIfNZ] = s
	IndentOpcodes[OpHandler] = s

	DedentOpcodes = make(map[Opcode]struct{})
	DedentOpcodes[OpEndDef] = s
	DedentOpcodes[OpElse] = s
	DedentOpcodes[OpEndIf] = s
}

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

// Dumper is a pointer to a function that may be passed to Run or Step, which
// if it is not nil, is called before each instruction. It can be used to dump
// the vm, record state, etc.
type Dumper func(*ChaincodeVM)

// Trace is a Dumper which emits the state of the stack before each instruction
func Trace(vm *ChaincodeVM) {
	fmt.Printf(
		"------\n%s\n%s\n",
		vm.Stack(),
		vm.DisassembleLine(vm.IP()),
	)
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

// MutableChaincodeVM is a chaincode vm into which instructions can be injected
// externally at arbitrary points.
//
// The default, normal sequence of operations is for a script to be immutable:
// it is created, initialized, run, and values are extracted. Determinism isn't
// just an accidental property; it's a designed requirement of the system.
//
// However, there are cases, such as within debuggers, in which it is desirable
// to be able to run arbitrary chaincode commands during the evaluation of a
// chaincode script. This breaks determinism for that context, but within the
// debugger, that isn't such a big problem.
//
// An instance of this type can be constructed from a normal ChaincodeVM
// by calling its MakeMutable method.
type MutableChaincodeVM struct {
	ChaincodeVM
}

// NewEmpty creates a new VM with a minimal empty handler
func NewEmpty() (*ChaincodeVM, error) {
	return New(ChasmBinary{"", "", MiniAsm("handler 0 enddef")})
}

// NewChaincode creates a new VM from the supplied chaincode
func NewChaincode(c Chaincode) (*ChaincodeVM, error) {
	return New(ChasmBinary{"", "", c})
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

// PreLoadOpcodes acepts an array of opcodes and validates it.
// If it fails to validate, the VM is not modified.
// However, if it does validate the VM is updated with
// code and function tables.
func (vm *ChaincodeVM) PreLoadOpcodes(data Chaincode) error {
	if err := data.IsValid(); err != nil {
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
func (vm *ChaincodeVM) Run(debug Dumper) error {
	if vm.runstate == RsReady {
		vm.runstate = RsRunning
	}
	for vm.runstate == RsRunning {
		if debug != nil {
			debug(vm)
		}
		if err := vm.Step(debug); err != nil {
			return err
		}
	}
	return nil
}

// Stringizer is used to override the default behavior of the
// disassembler for specific opcodes.
type Stringizer func(op Opcode, extra []Opcode) string

// DisasmHelpers is a map for specific opcodes to override the default renderer.
var DisasmHelpers = make(map[Opcode]Stringizer)

func init() {
	getInt := func(extra []Opcode) int64 {
		b := Chaincode(extra).Bytes()
		if len(b) < 8 {
			b = append(b, make([]byte, 8-len(b))...)
		}
		i := binary.LittleEndian.Uint64(b)
		return int64(i)
	}

	helpInt := func(op Opcode, extra []Opcode) string {
		return fmt.Sprintf("%-7s %-12d ", op, getInt(extra))
	}

	// add disassembly helpers which render the base-10 values
	// for all pushn opcodes
	DisasmHelpers[OpPush1] = helpInt
	DisasmHelpers[OpPush2] = helpInt
	DisasmHelpers[OpPush3] = helpInt
	DisasmHelpers[OpPush4] = helpInt
	DisasmHelpers[OpPush5] = helpInt
	DisasmHelpers[OpPush6] = helpInt
	DisasmHelpers[OpPush7] = helpInt
	DisasmHelpers[OpPush8] = helpInt

	// render a binary string
	DisasmHelpers[OpPushB] = func(op Opcode, extra []Opcode) string {
		data := fmt.Sprintf("%q", Chaincode(extra).Bytes())
		// trim surrounding quotes
		data = data[1 : len(data)-1]
		if len(data) > 21 { // empirically determined field width
			data = data[:18] + "..."
		}
		return fmt.Sprintf("%-7s %-12s ", op, data)
	}

	// render a timestamp
	DisasmHelpers[OpPushT] = func(op Opcode, extra []Opcode) string {
		return fmt.Sprintf("%-7s %-12s ", op, math.Timestamp(getInt(extra)).String())
	}
}

// DisassembleAll dumps a disassembly of the whole VM to the Writer
func (vm *ChaincodeVM) DisassembleAll(w io.Writer) {
	fmt.Fprintln(w, "--DISASSEMBLY--")
	tabs := 0
	for pc := 0; pc < len(vm.code); {
		// dedent before producing a line
		if _, ok := DedentOpcodes[vm.code[pc]]; ok {
			tabs--
		}

		// produce a line
		s, delta := vm.disassembleIndent(pc, tabs)
		fmt.Fprintln(w, s)

		// indent after producing a line
		if _, ok := IndentOpcodes[vm.code[pc]]; ok {
			tabs++
		}

		// go to next instruction
		pc += delta
	}
	fmt.Fprintln(w, "---------------")
}

// Disassemble returns a single disassembled instruction as a text string, possibly with embedded newlines,
// along with how many bytes it consumed.
func (vm *ChaincodeVM) Disassemble(pc int) (string, int) {
	return vm.disassembleIndent(pc, 0)
}

func (vm *ChaincodeVM) disassembleIndent(pc int, tabs int) (string, int) {
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
	out += fmt.Sprintf("%-24s  ", hex)

	for i := 0; i < tabs; i++ {
		out += tab
	}

	if helper, ok := DisasmHelpers[op]; !ok {
		args := ""
		if numExtra > 0 && numExtra < 5 {
			args = hex[3:]
		}
		if numExtra >= 5 {
			args = "..."
		}
		out += fmt.Sprintf("%-7s %-12s ", op, args)
	} else {
		out += helper(op, vm.code[pc+1:pc+1+numExtra])
	}

	return out, numExtra + 1
}

// String implements Stringer so we can print a VM and get something meaningful.
func (vm *ChaincodeVM) String() string {
	st := strings.Split(vm.stack.String(), "\n")
	st1 := make([]string, len(st))
	for i := range st {
		st1[len(st)-i-1] = st[i][4:]
	}
	disasm, _ := vm.Disassemble(vm.pc)
	return fmt.Sprintf("%-40s STK: %s\n", disasm, strings.Join(st1, ", "))
}

// Bytes returns the []byte corresponding to the chaincode
func (vm *ChaincodeVM) Bytes() []byte {
	b := make([]byte, len(vm.code))
	for i := range vm.code {
		b[i] = byte(vm.code[i])
	}
	return b
}

// DisassembledLine is the data structure intended to be leveraged by a
// debugging API.
type DisassembledLine struct {
	PC       int
	Opcode   Opcode
	NumExtra int
	ArgBytes []byte
}

func (l DisassembledLine) String() string {
	out := fmt.Sprintf("%04x: %s", l.PC, l.Opcode)
	for _, b := range l.ArgBytes {
		out += fmt.Sprintf(" %02x", b)
	}
	return out
}

// DisassembleLines returns a structured disassembly of the whole VM
// Do not call this on a vm that has not been validated!
func (vm *ChaincodeVM) DisassembleLines() []*DisassembledLine {
	var r []*DisassembledLine
	for pc := 0; pc < len(vm.code); {
		l := vm.DisassembleLine(pc)
		r = append(r, l)
		pc += l.NumExtra + 1
	}
	return r
}

// DisassembleLine returns a single disassembled instruction as an object
func (vm *ChaincodeVM) DisassembleLine(pc int) *DisassembledLine {
	if pc >= len(vm.code) {
		return nil
	}
	r := &DisassembledLine{
		PC:       pc,
		Opcode:   vm.code[pc],
		NumExtra: extraBytes(vm.code, pc),
	}
	if r.NumExtra > 0 {
		r.ArgBytes = make([]byte, r.NumExtra)
		for i := 1; i <= r.NumExtra; i++ {
			r.ArgBytes[i-1] = byte(vm.code[pc+i])
		}
	}

	return r
}

// MakeMutable allows this VM to be mutated externally.
//
// This function consumes the ChaincodeVM instance on which it is defined;
// no references to that instance should be retained or used after this is called.
// To help enforce this restriction, this function zeroizes the calling vm.
//
// This call is explicit and grep-able, so codebases which require determinism
// can easily prove that no script is mutable.
func (vm *ChaincodeVM) MakeMutable() *MutableChaincodeVM {
	mvm := MutableChaincodeVM{ChaincodeVM: *vm}
	*vm = ChaincodeVM{}
	return &mvm
}
