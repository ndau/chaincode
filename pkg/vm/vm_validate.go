package vm

import "fmt"

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
