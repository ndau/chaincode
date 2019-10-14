// Code generated automatically by "make generate"; DO NOT EDIT.

package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

// these are the opcodes supported by mini-asm
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
	"tuck":    OpTuck,
	"ret":     OpRet,
	"fail":    OpFail,
	"one":     OpOne,
	"neg1":    OpNeg1,
	"true":    OpTrue,
	"maxnum":  OpMaxNum,
	"minnum":  OpMinNum,
	"zero":    OpZero,
	"false":   OpFalse,
	"push1":   OpPush1,
	"push2":   OpPush2,
	"push3":   OpPush3,
	"push4":   OpPush4,
	"push5":   OpPush5,
	"push6":   OpPush6,
	"push7":   OpPush7,
	"push8":   OpPush8,
	"pushb":   OpPushB,
	"pusht":   OpPushT,
	"now":     OpNow,
	"rand":    OpRand,
	"pushl":   OpPushL,
	"add":     OpAdd,
	"sub":     OpSub,
	"mul":     OpMul,
	"div":     OpDiv,
	"mod":     OpMod,
	"divmod":  OpDivMod,
	"muldiv":  OpMulDiv,
	"not":     OpNot,
	"neg":     OpNeg,
	"inc":     OpInc,
	"dec":     OpDec,
	"index":   OpIndex,
	"len":     OpLen,
	"append":  OpAppend,
	"extend":  OpExtend,
	"slice":   OpSlice,
	"field":   OpField,
	"isfield": OpIsField,
	"fieldl":  OpFieldL,
	"def":     OpDef,
	"call":    OpCall,
	"deco":    OpDeco,
	"enddef":  OpEndDef,
	"ifz":     OpIfZ,
	"ifnz":    OpIfNZ,
	"else":    OpElse,
	"endif":   OpEndIf,
	"sum":     OpSum,
	"avg":     OpAvg,
	"max":     OpMax,
	"min":     OpMin,
	"choice":  OpChoice,
	"wchoice": OpWChoice,
	"sort":    OpSort,
	"lookup":  OpLookup,
	"handler": OpHandler,
	"or":      OpOr,
	"and":     OpAnd,
	"xor":     OpXor,
	"count1s": OpCount1s,
	"bnot":    OpBNot,
	"lt":      OpLt,
	"lte":     OpLte,
	"eq":      OpEq,
	"gte":     OpGte,
	"gt":      OpGt,
}
