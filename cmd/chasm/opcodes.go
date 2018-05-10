package main

// Opcodes
const (
	OpNop     byte = 0x00
	OpDrop    byte = 0x01
	OpDrop2   byte = 0x02
	OpDup     byte = 0x05
	OpDup2    byte = 0x06
	OpSwap    byte = 0x09
	OpOver    byte = 0x0D
	OpPick    byte = 0x0E
	OpRoll    byte = 0x0F
	OpRet     byte = 0x10
	OpFail    byte = 0x11
	OpZero    byte = 0x20
	OpFalse   byte = 0x20
	OpPushN   byte = 0x20
	OpPush64  byte = 0x29
	OpOne     byte = 0x2A
	OpTrue    byte = 0x2A
	OpNeg1    byte = 0x2B
	OpPushT   byte = 0x2C
	OpNow     byte = 0x2D
	OpRand    byte = 0x2F
	OpPushL   byte = 0x30
	OpAdd     byte = 0x40
	OpSub     byte = 0x41
	OpMul     byte = 0x42
	OpDiv     byte = 0x43
	OpMod     byte = 0x44
	OpNot     byte = 0x45
	OpNeg     byte = 0x46
	OpInc     byte = 0x47
	OpDec     byte = 0x48
	OpIndex   byte = 0x50
	OpLen     byte = 0x51
	OpAppend  byte = 0x52
	OpExtend  byte = 0x53
	OpSlice   byte = 0x54
	OpField   byte = 0x60
	OpFieldL  byte = 0x70
	OpIfz     byte = 0x80
	OpIfnz    byte = 0x81
	OpElse    byte = 0x82
	OpEnd     byte = 0x88
	OpSum     byte = 0x90
	OpAvg     byte = 0x91
	OpMax     byte = 0x92
	OpMin     byte = 0x93
	OpChoice  byte = 0x94
	OpWChoice byte = 0x95
	OpSort    byte = 0x96
	OpLookup  byte = 0x97
)
