package vm

import (
	"fmt"
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
	CtxAccount     ContextByte = iota
)

// Contexts is a map of ContextByte to context string
var Contexts = map[ContextByte]string{
	CtxTest:        "TEST",
	CtxNodePayout:  "NODE_PAYOUT",
	CtxEaiTiming:   "EAI_TIMING",
	CtxNodeQuality: "NODE_QUALITY",
	CtxMarketPrice: "MARKET_PRICE",
	CtxAccount:     "ACCOUNT",
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

// BuildVMForContext constructs a new VM in the desired context,
// checks to make sure that the desired context agrees with the VM's context,
// and populates the stack with whatever values are specified.
func BuildVMForContext(context ContextByte, bin ChasmBinary, values ...Value) (*ChaincodeVM, error) {
	vm, err := New(bin)
	if err != nil {
		return nil, err
	}
	if ContextByte(vm.context) != context {
		return nil, fmt.Errorf("binary context %d does not agree with required context %d", vm.context, context)
	}
	vm.Init(values...)
	return vm, nil
}

// BuildVMForTest constructs a new VM in the TEST context, and populates the
// stack with whatever values are specified.
func BuildVMForTest(bin ChasmBinary, values ...Value) (*ChaincodeVM, error) {
	return BuildVMForContext(CtxTest, bin, values...)
}

// We need parallel constructors for VMs that take the appropriate parameters
// on the build function for each type of Context, and construct a VM appropriately.
