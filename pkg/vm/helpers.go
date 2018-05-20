package vm

// Some masking values
const (
	ByteMask byte = 0xFF
	HighBit  byte = 0x80
)

// ToBytesU returns an array of 8 bytes encoding n as a uint in little-endian form
func ToBytesU(n uint64) []byte {
	b := []byte{}
	a := n
	for nbytes := 0; nbytes < 8; nbytes++ {
		b = append(b, byte(a)&ByteMask)
		a >>= 8
	}
	return b
}
