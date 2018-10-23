// Code generated automatically by "make generate" -- DO NOT EDIT

package vm

// extraBytes returns the number of extra bytes associated with a given opcode
func extraBytes(code Chaincode, offset int) int {
	// helper function for safety
	getat := func(ix int) Opcode {
		if ix >= len(code) {
			return 0
		}
		return code[ix]
	}

	numExtra := 0
	op := getat(offset)
	switch op {
	case OpPick:
		numExtra = 1
	case OpRoll:
		numExtra = 1
	case OpTuck:
		numExtra = 1
	case OpPush1:
		numExtra = 1
	case OpPush2:
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
	case OpPush8:
		numExtra = 8
	case OpPushB:
		numExtra = int(getat(offset+1)) + 1
	case OpPushT:
		numExtra = 8
	case OpPushA:
		numExtra = int(getat(offset+1)) + 1
	case OpField:
		numExtra = 1
	case OpIsField:
		numExtra = 1
	case OpFieldL:
		numExtra = 1
	case OpDef:
		numExtra = 2
	case OpCall:
		numExtra = 1
	case OpDeco:
		numExtra = 2
	case OpWChoice:
		numExtra = 1
	case OpSort:
		numExtra = 1
	case OpLookup:
		numExtra = 1
	case OpHandler:
		numExtra = int(getat(offset+1)) + 1
	}
	return numExtra
}
