package vm

import (
	"strings"
)

// ContextByte is a byte used to identify context
type ContextByte byte

// Constants for Contexts
const (
	CtxTest        ContextByte = iota
	CtxNodePayout  ContextByte = iota
	CtxEaiTiming   ContextByte = iota
	CtxNodeQuality ContextByte = iota
	CtxMarketPrice ContextByte = iota
)

// Contexts is a map of ContextByte to context string
var Contexts = map[ContextByte]string{
	CtxTest:        "TEST",
	CtxNodePayout:  "NODE_PAYOUT",
	CtxEaiTiming:   "EAI_TIMING",
	CtxNodeQuality: "NODE_QUALITY",
	CtxMarketPrice: "MARKET_PRICE",
}

// ContextLookup searches for a context by a given name; returns true if found
func ContextLookup(s string) (ContextByte, bool) {
	for k, v := range Contexts {
		if strings.EqualFold(v, s) {
			return k, true
		}
	}
	return CtxTest, false
}
