package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

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
	StHandler structureState = iota
	StHndPlus structureState = iota
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
func validateStructure(code Chaincode) (map[byte]int, []funcInfo, error) {
	// This is a table of allowed transitions from a given state, depending on the opcode.
	// For example, from the StNull state, the only allowed transitions are to definitions.
	transitions := map[tr]structureState{
		tr{StNull, OpDef}:        StDefPlus,
		tr{StNull, OpHandler}:    StHndPlus,
		tr{StNull, OpIfZ}:        StError,
		tr{StNull, OpIfNZ}:       StError,
		tr{StNull, OpElse}:       StError,
		tr{StNull, OpEndIf}:      StError,
		tr{StNull, OpEndDef}:     StError,
		tr{StDef, OpDef}:         StError,
		tr{StDef, OpHandler}:     StError,
		tr{StDef, OpIfZ}:         StIfPlus,
		tr{StDef, OpIfNZ}:        StIfPlus,
		tr{StDef, OpElse}:        StError,
		tr{StDef, OpEndIf}:       StError,
		tr{StDef, OpEndDef}:      StNull,
		tr{StHandler, OpDef}:     StError,
		tr{StHandler, OpHandler}: StError,
		tr{StHandler, OpIfZ}:     StIfPlus,
		tr{StHandler, OpIfNZ}:    StIfPlus,
		tr{StHandler, OpElse}:    StError,
		tr{StHandler, OpEndIf}:   StError,
		tr{StHandler, OpEndDef}:  StNull,
		tr{StIf, OpDef}:          StError,
		tr{StIf, OpHandler}:      StError,
		tr{StIf, OpIfZ}:          StIfPlus,
		tr{StIf, OpIfNZ}:         StIfPlus,
		tr{StIf, OpElse}:         StElse,
		tr{StIf, OpEndIf}:        StIfMinus,
		tr{StIf, OpEndDef}:       StError,
		tr{StElse, OpDef}:        StError,
		tr{StElse, OpHandler}:    StError,
		tr{StElse, OpIfZ}:        StIfPlus,
		tr{StElse, OpIfNZ}:       StIfPlus,
		tr{StElse, OpElse}:       StError,
		tr{StElse, OpEndIf}:      StIfMinus,
		tr{StElse, OpEndDef}:     StError,
	}

	var state structureState
	var depth int
	var nfuncs byte
	var skipcount int
	var handlers = make(map[byte]int)
	var functions = make([]funcInfo, 0)

	// for offset, b := range code {
	for offset := 0; offset < len(code); offset += skipcount + 1 {
		// calculate how many extra bytes we need for this opcode
		skipcount = extraBytes(code, offset)
		// if we don't have that many left in the code block, it's bad
		if offset+skipcount >= len(code) {
			return handlers, functions, ValidationError{"missing operands"}
		}

		newstate, found := transitions[tr{state, code[offset]}]
		if !found {
			continue
		}
		switch newstate {
		case StDefPlus:
			// this is the old code that should be preserved with offsets renamed to functions
			funcnum := byte(code[offset+1])
			if funcnum != nfuncs {
				return handlers, functions, ValidationError{fmt.Sprintf("def should have been %d, found %d", nfuncs, funcnum)}
			}
			functions = append(functions, funcInfo{
				offset: int(offset + skipcount + 1), // skip the def opcode and parms
				nargs:  int(code[offset+2]),
			})
			nfuncs++
			newstate = StDef
		case StHndPlus:
			nhandlers := int(code[offset+1])
			if len(code) < offset+2+nhandlers {
				return handlers, functions, ValidationError{"handler count parameter was too large"}
			}
			handlerids := code[offset+2 : offset+2+nhandlers]
			if nhandlers == 0 {
				// special case -- "handler 0" means define the default
				handlerids = []Opcode{0}
			}
			for i := 0; i < len(handlerids); i++ {
				handlerID := byte(handlerids[i])
				if _, found := handlers[handlerID]; found {
					return handlers, functions, ValidationError{fmt.Sprintf("multiple handlers found for event %d", handlerID)}
				}
				handlers[handlerID] = int(offset + skipcount + 1) // skip the handler opcode
			}
			newstate = StHandler
		case StError:
			return handlers, functions, ValidationError{fmt.Sprintf("invalid program structure [offset=%d]", offset)}
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
		return handlers, functions, ValidationError{"missing operands"}
	}
	if state != StNull {
		return handlers, functions, ValidationError{"missing end"}
	}
	if len(handlers) == 0 {
		return handlers, functions, ValidationError{"no handlers were defined"}
	}
	return handlers, functions, nil
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
		bitset.Set(byte(instrs[i][0]))
	}
	return bitset
}

// bitsetToOpcodes returns a human-readable list of opcodes as a bitset
func bitsetToOpcodes(b *bitset256.Bitset256) string {
	sa := []string{}
	for i := 0; i < 256; i++ {
		if b.Get(byte(i)) {
			sa = append(sa, Opcode(i).String())
		}
	}
	return strings.Join(sa, " ")
}

// DisableOpcode allows an opcode to be disabled at runtime; this is for use in
// a security situation where a vulnerability has been discovered and it's
// important to make sure that VMs containing this opcode can no longer run.
//
// Note that this has global impact and cannot be reversed! This is to say that
// there is no equivalent "EnableOpcode" function. Once an opcode is disabled,
// the only way to re-enable it is to restart the application.
//
// This operates at the level of VM validation -- a VM that is already loaded
// will not be affected by this operation. The function returns true if the
// opcode was previously enabled (i.e., if it has had an effect).
func DisableOpcode(op Opcode) bool {
	ret := EnabledOpcodes.Get(byte(op))
	EnabledOpcodes.Clear(byte(op))
	return ret
}

// CodeSize calculates the size of code without counting the size of the
// data included in PushB opcodes. It should only be called
// after validateStructure.
func (code Chaincode) CodeSize() int {
	size := 0
	skipcount := 0
	for offset := 0; offset < len(code); offset += skipcount + 1 {
		skipcount = extraBytes(code, offset)
		switch code[offset] {
		case OpPushB:
			// don't count the data bytes
			size += 2
		default:
			size += skipcount + 1
		}
	}
	return size
}
