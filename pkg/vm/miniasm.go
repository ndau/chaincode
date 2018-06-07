package vm

import (
	"regexp"
	"strconv"
	"strings"
)

// MiniAsm is a miniature assembler that has very simple syntax. It's primarily intended for writing
// simple test code.
//
// It takes a single string as input. It then converts everything to lower case and splits it into
// 'words' by whitespace -- all whitespace is equivalent.
//
// If a word matches an opcode, it generates the associated opcode.
// If a word matches a simple pattern for a timestamp, it attempts to parse it as a timestamp.
// All other words are expected to be one-byte hex values.
//
// Any failure in parsing causes a panic; there is no error recovery.
//
// The resulting stream of instructions is returned prefixed with a 0 byte, which is the "TEST" context.
// No attempt is made to ensure that opcode parameters or types are correct, and each opcode is individually
// specified (Push1, Push2, etc).
//

var opcodeMap = map[string]Opcode{
	"nop":     OpNop,
	"drop":    OpDrop,
	"drop2":   OpDrop2,
	"dup":     OpDup,
	"dup2":    OpDup2,
	"swap":    OpSwap,
	"over":    OpOver,
	"pick":    OpPick,
	"roll":    OpRoll,
	"ret":     OpRet,
	"fail":    OpFail,
	"zero":    OpZero,
	"false":   OpFalse,
	"pushN":   OpPushN,
	"push1":   OpPush1,
	"push2":   OpPush2,
	"push3":   OpPush3,
	"push4":   OpPush4,
	"push5":   OpPush5,
	"push6":   OpPush6,
	"push7":   OpPush7,
	"push8":   OpPush8,
	"pushb":   OpPushB,
	"one":     OpOne,
	"true":    OpTrue,
	"neg1":    OpNeg1,
	"pusht":   OpPushT,
	"now":     OpNow,
	"rand":    OpRand,
	"pushl":   OpPushL,
	"add":     OpAdd,
	"sub":     OpSub,
	"mul":     OpMul,
	"div":     OpDiv,
	"mod":     OpMod,
	"not":     OpNot,
	"neg":     OpNeg,
	"inc":     OpInc,
	"dec":     OpDec,
	"eq":      OpEq,
	"lt":      OpLt,
	"gt":      OpGt,
	"index":   OpIndex,
	"len":     OpLen,
	"append":  OpAppend,
	"extend":  OpExtend,
	"slice":   OpSlice,
	"field":   OpField,
	"fieldl":  OpFieldL,
	"def":     OpDef,
	"call":    OpCall,
	"deco":    OpDeco,
	"ifz":     OpIfz,
	"ifnz":    OpIfnz,
	"else":    OpElse,
	"enddef":  OpEndDef,
	"endif":   OpEndIf,
	"sum":     OpSum,
	"avg":     OpAvg,
	"max":     OpMax,
	"min":     OpMin,
	"choice":  OpChoice,
	"wchoice": OpWChoice,
	"sort":    OpSort,
	"lookup":  OpLookup,
}

// miniAsm is primarily for testing but we want it available.
// nolint: deadcode
func miniAsm(s string) []Opcode {
	// whitespace
	wsp := regexp.MustCompile("[ \t\r\n]")
	// timestamp
	tsp := regexp.MustCompile("[0-9-]+T[0-9:]+Z")
	// quoted string without spaces (this is a mini assembler!)
	qsp := regexp.MustCompile(`"[^" ]+"`)
	words := wsp.Split(strings.TrimSpace(s), -1)
	opcodes := []Opcode{0}
	for _, w := range words {
		// skip empty words
		if w == "" {
			continue
		}
		// see if it's an opcode
		if op, ok := opcodeMap[strings.ToLower(w)]; ok {
			opcodes = append(opcodes, op)
			continue
		}
		// see if it's a timestamp
		if tsp.MatchString(strings.ToUpper(w)) {
			t, err := ParseTimestamp(strings.ToUpper(w))
			if err != nil {
				panic(err)
			}
			bytes := ToBytes(int64(t.t))
			for _, byt := range bytes {
				opcodes = append(opcodes, Opcode(byt))
			}
			continue
		}
		// see if it's a quoted string
		if qsp.MatchString(w) {
			bytes := w[1 : len(w)-1]
			opcodes = append(opcodes, Opcode(len(bytes)))
			for _, byt := range bytes {
				opcodes = append(opcodes, Opcode(byt))
			}
			continue
		}
		// otherwise it should be a hex value
		b, err := strconv.ParseUint(w, 16, 8)
		if err != nil {
			panic(err)
		}
		opcodes = append(opcodes, Opcode(b))
	}
	return opcodes
}
