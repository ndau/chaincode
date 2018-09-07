// This file is generated automatically; DO NOT EDIT.

package vm

// Opcode is a byte used to identify an opcode; we rely on it being a byte in some cases.
type Opcode byte

//go:generate stringer -trimprefix Op -type Opcode opcodes.go
//go:generate msgp -tests=0

// Opcodes
const (
	OpNop     Opcode = 0x00
	OpDrop    Opcode = 0x01
	OpDrop2   Opcode = 0x02
	OpDup     Opcode = 0x05
	OpDup2    Opcode = 0x06
	OpSwap    Opcode = 0x09
	OpOver    Opcode = 0x0c
	OpPick    Opcode = 0x0d
	OpRoll    Opcode = 0x0e
	OpTuck    Opcode = 0x0f
	OpRet     Opcode = 0x10
	OpFail    Opcode = 0x11
	OpOne     Opcode = 0x1a
	OpNeg1    Opcode = 0x1b
	OpTrue    Opcode = 0x1b
	OpMaxNum  Opcode = 0x1c
	OpMinNum  Opcode = 0x1d
	OpZero    Opcode = 0x20
	OpFalse   Opcode = 0x20
	OpPush1   Opcode = 0x21
	OpPush2   Opcode = 0x22
	OpPush3   Opcode = 0x23
	OpPush4   Opcode = 0x24
	OpPush5   Opcode = 0x25
	OpPush6   Opcode = 0x26
	OpPush7   Opcode = 0x27
	OpPush8   Opcode = 0x28
	OpPushB   Opcode = 0x2a
	OpPushT   Opcode = 0x2b
	OpNow     Opcode = 0x2c
	OpPushA   Opcode = 0x2d
	OpRand    Opcode = 0x2e
	OpPushL   Opcode = 0x2f
	OpAdd     Opcode = 0x40
	OpSub     Opcode = 0x41
	OpMul     Opcode = 0x42
	OpDiv     Opcode = 0x43
	OpMod     Opcode = 0x44
	OpDivMod  Opcode = 0x45
	OpMulDiv  Opcode = 0x46
	OpNot     Opcode = 0x48
	OpNeg     Opcode = 0x49
	OpInc     Opcode = 0x4a
	OpDec     Opcode = 0x4b
	OpIndex   Opcode = 0x50
	OpLen     Opcode = 0x51
	OpAppend  Opcode = 0x52
	OpExtend  Opcode = 0x53
	OpSlice   Opcode = 0x54
	OpField   Opcode = 0x60
	OpIsField Opcode = 0x61
	OpFieldL  Opcode = 0x70
	OpDef     Opcode = 0x80
	OpCall    Opcode = 0x81
	OpDeco    Opcode = 0x82
	OpEndDef  Opcode = 0x88
	OpIfZ     Opcode = 0x89
	OpIfNZ    Opcode = 0x8a
	OpElse    Opcode = 0x8e
	OpEndIf   Opcode = 0x8f
	OpSum     Opcode = 0x90
	OpAvg     Opcode = 0x91
	OpMax     Opcode = 0x92
	OpMin     Opcode = 0x93
	OpChoice  Opcode = 0x94
	OpWChoice Opcode = 0x95
	OpSort    Opcode = 0x96
	OpLookup  Opcode = 0x97
	OpHandler Opcode = 0xa0
	OpOr      Opcode = 0xb0
	OpAnd     Opcode = 0xb1
	OpXor     Opcode = 0xb2
	OpCount1s Opcode = 0xbc
	OpBNot    Opcode = 0xbf
	OpLt      Opcode = 0xc0
	OpLte     Opcode = 0xc1
	OpEq      Opcode = 0xc2
	OpGte     Opcode = 0xc3
	OpGt      Opcode = 0xc4
)
