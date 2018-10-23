// Code generated automatically by "make generate" -- DO NOT EDIT.

package vm

import "github.com/oneiro-ndev/ndaumath/pkg/bitset256"

// EnabledOpcodes is a bitset of the opcodes that are enabled -- only these opcodes will be
// permitted in scripts.
var EnabledOpcodes = bitset256.New(
	byte(OpNop),
	byte(OpDrop),
	byte(OpDrop2),
	byte(OpDup),
	byte(OpDup2),
	byte(OpSwap),
	byte(OpOver),
	byte(OpPick),
	byte(OpRoll),
	byte(OpTuck),
	byte(OpRet),
	byte(OpFail),
	byte(OpOne),
	byte(OpNeg1),
	byte(OpMaxNum),
	byte(OpMinNum),
	byte(OpZero),
	byte(OpPush1),
	byte(OpPush2),
	byte(OpPush3),
	byte(OpPush4),
	byte(OpPush5),
	byte(OpPush6),
	byte(OpPush7),
	byte(OpPush8),
	byte(OpPushB),
	byte(OpPushT),
	byte(OpNow),
	byte(OpRand),
	byte(OpPushL),
	byte(OpAdd),
	byte(OpSub),
	byte(OpMul),
	byte(OpDiv),
	byte(OpMod),
	byte(OpDivMod),
	byte(OpMulDiv),
	byte(OpNot),
	byte(OpNeg),
	byte(OpInc),
	byte(OpDec),
	byte(OpIndex),
	byte(OpLen),
	byte(OpAppend),
	byte(OpExtend),
	byte(OpSlice),
	byte(OpField),
	byte(OpIsField),
	byte(OpFieldL),
	byte(OpDef),
	byte(OpCall),
	byte(OpDeco),
	byte(OpEndDef),
	byte(OpIfZ),
	byte(OpIfNZ),
	byte(OpElse),
	byte(OpEndIf),
	byte(OpSum),
	byte(OpAvg),
	byte(OpMax),
	byte(OpMin),
	byte(OpChoice),
	byte(OpWChoice),
	byte(OpSort),
	byte(OpLookup),
	byte(OpHandler),
	byte(OpOr),
	byte(OpAnd),
	byte(OpXor),
	byte(OpCount1s),
	byte(OpBNot),
	byte(OpLt),
	byte(OpLte),
	byte(OpEq),
	byte(OpGte),
	byte(OpGt),
)
