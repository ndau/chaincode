package vm

// Constants for Contexts
const (
	CtxTest        byte = iota
	CtxNodePayout  byte = iota
	CtxEaiTiming   byte = iota
	CtxNodeQuality byte = iota
	CtxMarketPrice byte = iota
)

// Contexts is a map of context byte to context string
var Contexts = map[byte]string{
	CtxTest:        "TEST",
	CtxNodePayout:  "NODE_PAYOUT",
	CtxEaiTiming:   "EAI_TIMING",
	CtxNodeQuality: "NODE_QUALITY",
	CtxMarketPrice: "MARKET_PRICE",
}
