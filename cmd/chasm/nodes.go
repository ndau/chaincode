package main

import (
	"errors"
	"strconv"

	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

// Node is the fundamental unit that the parser manipulates (it builds a structure of nodes).
//Each node can emit itself as an array of bytes, or nil.
type Node interface {
	bytes() []byte
}

// Fixupper is an interface that is implemented by all nodes that need fixups and all nodes
// that contain other nodes as children. It is called before the bytes() function to allow
// nodes to do any fixing up necessary.
type Fixupper interface {
	fixup(map[string]int)
}

// Script is the highest level node in the system
type Script struct {
	preamble Node
	opcodes  []Node
}

func (n *Script) fixup(funcs map[string]int) {
	for _, op := range n.opcodes {
		if f, ok := op.(Fixupper); ok {
			f.fixup(funcs)
		}
	}
}

func (n *Script) bytes() []byte {
	b := append([]byte{}, n.preamble.bytes()...)
	for _, op := range n.opcodes {
		b = append(b, op.bytes()...)
	}
	return b
}

func newScript(p interface{}, opcodes interface{}) (*Script, error) {
	preamble, ok := p.(*PreambleNode)
	if !ok {
		return &Script{}, errors.New("not a preamble node")
	}
	sl := toIfaceSlice(opcodes)
	ops := []Node{}
	for _, v := range sl {
		if n, ok := v.(Node); ok {
			ops = append(ops, n)
		}
	}
	return &Script{preamble: preamble, opcodes: ops}, nil
}

// PreambleNode expresses the information in the preamble (which for now is just a context byte)
type PreambleNode struct {
	context vm.ContextByte
}

func (n *PreambleNode) bytes() []byte {
	return []byte{byte(n.context)}
}

func newPreambleNode(ctx vm.ContextByte) (*PreambleNode, error) {
	return &PreambleNode{context: ctx}, nil
}

// FunctionDef is a node that expresses the information in a function definition
type FunctionDef struct {
	name  string
	index byte
	nodes []Node
}

func (n *FunctionDef) fixup(funcs map[string]int) {
	me, ok := funcs[n.name]
	if ok {
		n.index = byte(me)
	}
	for _, op := range n.nodes {
		if f, ok := op.(Fixupper); ok {
			f.fixup(funcs)
		}
	}
}

func (n *FunctionDef) bytes() []byte {
	b := []byte{byte(vm.OpDef), byte(n.index)}
	for _, op := range n.nodes {
		b = append(b, op.bytes()...)
	}
	b = append(b, byte(vm.OpEndDef))
	return b
}

func newFunctionDef(name string, nodes interface{}) (*FunctionDef, error) {
	sl := toIfaceSlice(nodes)
	nl := []Node{}
	for _, v := range sl {
		if n, ok := v.(Node); ok {
			nl = append(nl, n)
		}
	}
	f := &FunctionDef{name: name, index: 0xff, nodes: nl}
	return f, nil
}

// UnitaryOpcode is for opcodes that cannot take arguments
type UnitaryOpcode struct {
	opcode vm.Opcode
}

func (n *UnitaryOpcode) bytes() []byte {
	return []byte{byte(n.opcode)}
}

func newUnitaryOpcode(op vm.Opcode) (*UnitaryOpcode, error) {
	return &UnitaryOpcode{opcode: op}, nil
}

// BinaryOpcode is for opcodes that take one single-byte argument
type BinaryOpcode struct {
	opcode vm.Opcode
	value  byte
}

func (n BinaryOpcode) bytes() []byte {
	return []byte{byte(n.opcode), n.value}
}

func newBinaryOpcode(op vm.Opcode, v string) (*BinaryOpcode, error) {
	n, err := strconv.ParseUint(v, 0, 8)
	if err != nil {
		return &BinaryOpcode{}, err
	}
	return &BinaryOpcode{opcode: op, value: byte(n)}, nil
}

// CallOpcode is for opcodes that call a function and take a function name
type CallOpcode struct {
	opcode vm.Opcode
	name   string
	fix    byte
	value  byte
}

func (n *CallOpcode) fixup(funcs map[string]int) {
	me, ok := funcs[n.name]
	if ok {
		n.fix = byte(me)
	}
}

func (n *CallOpcode) bytes() []byte {
	return []byte{byte(n.opcode), n.fix, n.value}
}

func newCallOpcode(op vm.Opcode, name string, v string) (*CallOpcode, error) {
	n, err := strconv.ParseUint(v, 0, 8)
	if err != nil {
		return &CallOpcode{}, err
	}
	return &CallOpcode{opcode: op, name: name, value: byte(n)}, nil
}

// PushOpcode constructs push operations with the appropriate number of bytes to express
// the specified value. It has special cases for the special opcodes zero, one, and neg1.
type PushOpcode struct {
	arg int64
}

// This function builds a sequence of bytes consisting of either:
//   A ZERO, ONE, or NEG1 opcode
// OR
//   A PushN opcode followed by N bytes, where N is a value from 1-8.
//   The bytes are a representation of the value in little-endian order (low
//   byte first). The highest bit is the sign bit.
func (n *PushOpcode) bytes() []byte {
	switch n.arg {
	case 0:
		return []byte{byte(vm.OpZero)}
	case 1:
		return []byte{byte(vm.OpOne)}
	case -1:
		return []byte{byte(vm.OpNeg1)}
	default:
		b := vm.ToBytes(n.arg)
		var suppress byte
		if n.arg < 0 {
			suppress = byte(0xFF)
		}
		for b[len(b)-1] == suppress {
			b = b[:len(b)-1]
		}
		nbytes := byte(len(b))
		op := byte(vm.OpPushN) | nbytes
		b = append([]byte{op}, b...)
		return b
	}
}

func newPushOpcode(s string) (*PushOpcode, error) {
	v, err := strconv.ParseInt(s, 0, 64)
	return &PushOpcode{arg: v}, err
}

// Push64 is a 64-bit unsigned value
type Push64 struct {
	u uint64
}

func newPush64(s string) (*Push64, error) {
	v, err := strconv.ParseUint(s, 0, 64)
	return &Push64{u: v}, err
}

func (n *Push64) bytes() []byte {
	return append([]byte{byte(vm.OpPush64)}, vm.ToBytesU(n.u)...)
}

// PushTimestamp is a 64-bit representation of the time since the start of the epoch in microseconds
type PushTimestamp struct {
	t uint64
}

func newPushTimestamp(s string) (*PushTimestamp, error) {
	ts, err := vm.ParseTimestamp(s)
	if err != nil {
		return &PushTimestamp{}, err
	}
	return &PushTimestamp{ts.T()}, nil
}

func (n *PushTimestamp) bytes() []byte {
	return append([]byte{byte(vm.OpPushT)}, vm.ToBytesU(n.t)...)
}
