package vm

// Opcode is a byte used to identify an opcode; we rely on it being a byte in some cases.
type Opcode byte

//go:generate stringer -trimprefix Op -type Opcode opcodes.go

// Opcodes
const (
	OpNop     Opcode = 0x00
	OpDrop    Opcode = 0x01
	OpDrop2   Opcode = 0x02
	OpDup     Opcode = 0x05
	OpDup2    Opcode = 0x06
	OpSwap    Opcode = 0x09
	OpOver    Opcode = 0x0D
	OpPick    Opcode = 0x0E
	OpRoll    Opcode = 0x0F
	OpRet     Opcode = 0x10
	OpFail    Opcode = 0x11
	OpZero    Opcode = 0x20
	OpFalse   Opcode = 0x20
	OpPushN   Opcode = 0x20
	OpPush1   Opcode = 0x21
	OpPush2   Opcode = 0x22
	OpPush3   Opcode = 0x23
	OpPush4   Opcode = 0x24
	OpPush5   Opcode = 0x25
	OpPush6   Opcode = 0x26
	OpPush7   Opcode = 0x27
	OpPush8   Opcode = 0x28
	OpPush64  Opcode = 0x29
	OpOne     Opcode = 0x2A
	OpTrue    Opcode = 0x2A
	OpNeg1    Opcode = 0x2B
	OpPushT   Opcode = 0x2C
	OpNow     Opcode = 0x2D
	OpRand    Opcode = 0x2F
	OpPushL   Opcode = 0x30
	OpAdd     Opcode = 0x40
	OpSub     Opcode = 0x41
	OpMul     Opcode = 0x42
	OpDiv     Opcode = 0x43
	OpMod     Opcode = 0x44
	OpNot     Opcode = 0x45
	OpNeg     Opcode = 0x46
	OpInc     Opcode = 0x47
	OpDec     Opcode = 0x48
	OpEq      Opcode = 0x49
	OpGt      Opcode = 0x4A
	OpLt      Opcode = 0x4B
	OpIndex   Opcode = 0x50
	OpLen     Opcode = 0x51
	OpAppend  Opcode = 0x52
	OpExtend  Opcode = 0x53
	OpSlice   Opcode = 0x54
	OpField   Opcode = 0x60
	OpFieldL  Opcode = 0x70
	OpDef     Opcode = 0x80
	OpCall    Opcode = 0x81
	OpDeco    Opcode = 0x82
	OpEndDef  Opcode = 0x88
	OpIfz     Opcode = 0x89
	OpIfnz    Opcode = 0x8A
	OpElse    Opcode = 0x8E
	OpEndIf   Opcode = 0x8F
	OpSum     Opcode = 0x90
	OpAvg     Opcode = 0x91
	OpMax     Opcode = 0x92
	OpMin     Opcode = 0x93
	OpChoice  Opcode = 0x94
	OpWChoice Opcode = 0x95
	OpSort    Opcode = 0x96
	OpLookup  Opcode = 0x97
)
