package vm

import (
	"errors"
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
	CtxTransaction ContextByte = iota
)

// Contexts is a map of ContextByte to context string
var Contexts = map[ContextByte]string{
	CtxTest:        "TEST",
	CtxNodePayout:  "NODE_PAYOUT",
	CtxEaiTiming:   "EAI_TIMING",
	CtxNodeQuality: "NODE_QUALITY",
	CtxMarketPrice: "MARKET_PRICE",
	CtxTransaction: "TRANSACTION",
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

func BuildVmForTest(bin ChasmBinary) (*ChaincodeVM, error) {
	vm, err := New(bin)
	if err != nil {
		return nil, err
	}
	if ContextByte(vm.context) != CtxTest {
		return nil, errors.New("binary does not have required context")
	}
	// Test context has no initial stack
	vm.Init()
	return vm, nil
}

// We need parallel constructors for VMs that take the appropriate parameters
// on the build function for each type of Context, and construct a VM appropriately.
