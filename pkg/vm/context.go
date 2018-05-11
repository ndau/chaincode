package vm

import (
	"strings"
)

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

// ContextLookup searches for a context by a given name; returns true if found
func ContextLookup(s string) (byte, bool) {
	for k, v := range Contexts {
		if strings.EqualFold(v, s) {
			return k, true
		}
	}
	return CtxTest, false
}
