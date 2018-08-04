package vm

import (
	"fmt"
	"strings"

	"github.com/oneiro-ndev/ndaumath/pkg/bitset256"
)

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
		tr{StNull, OpIfZ}:    StError,
		tr{StNull, OpIfNZ}:   StError,
		tr{StNull, OpElse}:   StError,
		tr{StNull, OpEndIf}:  StError,
		tr{StNull, OpEndDef}: StError,
		tr{StDef, OpDef}:     StError,
		tr{StDef, OpIfZ}:     StIfPlus,
		tr{StDef, OpIfNZ}:    StIfPlus,
		tr{StDef, OpElse}:    StError,
		tr{StDef, OpEndIf}:   StError,
		tr{StDef, OpEndDef}:  StNull,
		tr{StIf, OpDef}:      StError,
		tr{StIf, OpIfZ}:      StIfPlus,
		tr{StIf, OpIfNZ}:     StIfPlus,
		tr{StIf, OpElse}:     StElse,
		tr{StIf, OpEndIf}:    StIfMinus,
		tr{StIf, OpEndDef}:   StError,
		tr{StElse, OpDef}:    StError,
		tr{StElse, OpIfZ}:    StIfPlus,
		tr{StElse, OpIfNZ}:   StIfPlus,
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

// generateInstructions is a helper that reads a sequence of bytes that has
// already been determined to be structurally valid and converts it to a set of Instruction
// objects, each of which is a slice consisting of a single opcode plus
// its data bytes
func generateInstructions(code []Opcode) []Instruction {
	instrs := make([]Instruction, 0)
	for pc := 0; pc < len(code); {
		// op := code[pc]
		numExtra := extraBytes(code, pc)
		inst := code[pc : pc+numExtra+1]
		instrs = append(instrs, inst)
		pc += numExtra + 1
	}
	return instrs
}

// getUsedOpcodes takes an array of Instructions and pulls out a bitset of the
// opcodes that were used in it.
func getUsedOpcodes(instrs []Instruction) *bitset256.Bitset256 {
	bitset := bitset256.New()
	for i := 0; i < len(instrs); i++ {
		bitset.Set(int(instrs[i][0]))
	}
	return bitset
}

// bitsetToOpcodes returns a human-readable list of opcodes as a bitset
func bitsetToOpcodes(b *bitset256.Bitset256) string {
	sa := []string{}
	for i := 0; i < 256; i++ {
		if b.Get(i) {
			sa = append(sa, Opcode(i).String())
		}
	}
	return strings.Join(sa, " ")
}

// DisableOpcode allows an opcode to be disabled at runtime; this is for use
// in a security situation where a vulnerability has been discovered and
// it's important to make sure that VMs containing this opcode can no longer
// run.
// Note that this has global impact and cannot be reversed! Once an opcode
// is disabled, the only way to re-enable it is to restart the application.
// Note that this operates at the level of VM validation -- a VM that is
// already loaded will not be affected by this operation.
// The function returns true if the opcode was previously enabled (i.e., if it
// has had an effect).
func DisableOpcode(op Opcode) bool {
	ret := EnabledOpcodes.Get(int(op))
	EnabledOpcodes.Clear(int(op))
	return ret
}
