package vm

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/oneiro-ndev/ndaumath/pkg/signed"
)

func (vm *ChaincodeVM) runtimeError(err error) error {
	rte := wrapRuntimeError(err)
	return rte.PC(vm.pc - 1)
}

// This is only run on VMs that have been validated
func (vm *ChaincodeVM) skipToMatchingBracket() error {
	for {
		instr := vm.code[vm.pc]
		extra := extraBytes(vm.code, vm.pc)
		vm.pc += extra + 1
		nesting := 0
		switch instr {
		case OpIfNZ, OpIfZ:
			nesting++
		case OpElse:
			if nesting == 0 {
				// we're at the right level, so we're done
				return nil
			}
		case OpEndIf:
			if nesting > 0 {
				nesting--
			} else {
				// we're at the right level so we're done
				return nil
			}
		default:
			if vm.pc > len(vm.code) {
				// fail-safe (should never happen)
				panic("VM RAN OFF THE END!")
			}
		}
	}
}

// callFunction calls the function numbered by funcnum, copying nargs to a new stack
// returns the value left on the stack by the called function
func (vm *ChaincodeVM) callFunction(funcnum int, nargs int, debug bool, extraArgs ...Value) (Value, error) {
	var retval Value
	if funcnum <= vm.infunc || funcnum >= len(vm.functions) {
		return retval, vm.runtimeError(newRuntimeError("invalid function number (no recursion allowed)"))
	}
	newpc := vm.functions[funcnum]

	childvm, err := vm.CreateForFunc(funcnum, newpc, nargs)
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	for _, e := range extraArgs {
		err := childvm.stack.Push(e)
		if err != nil {
			return retval, vm.runtimeError(err)
		}
	}
	err = childvm.Run(debug)
	// no matter what, we want the history
	vm.history = append(vm.history, childvm.history...)
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	// we've called the child function, now get back its return value
	retval, err = childvm.stack.Pop()
	if err != nil {
		return retval, vm.runtimeError(err)
	}
	return retval, nil
}

// Step executes a single instruction
func (vm *ChaincodeVM) Step(debug bool) error {
	switch vm.runstate {
	default:
		return newRuntimeError("vm is not in runnable state")
	case RsReady:
		vm.runstate = RsRunning
		fallthrough
	case RsRunning:
		vm.history = append(vm.history, HistoryState{
			PC:    vm.pc,
			Stack: vm.stack.Clone(),
			// lists: vm.lists[:],
		})
	}

	// Check to see if we're still in runnable code
	if vm.pc >= len(vm.code) {
		vm.runstate = RsComplete
		return nil
	}

	// fetch the instruction
	instr := vm.code[vm.pc]
	// we always increment the pc immediately; we may also add to it if an instruction has additional data
	// when we report errors, we subtract 1 to get the correct value
	vm.pc++
	switch instr {
	case OpNop:
		// do nothing
	case OpDrop:
		// discards top element on stack
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
	case OpDrop2:
		// discards top two elements on stack
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
		if _, err := vm.stack.Pop(); err != nil {
			return vm.runtimeError(err)
		}
	case OpDup:
		v, err := vm.stack.Peek()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpDup2:
		v1, err := vm.stack.Get(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		v0, err := vm.stack.Get(0)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v0); err != nil {
			return vm.runtimeError(err)
		}
	case OpSwap:
		v1, err := vm.stack.PopAt(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
	case OpOver:
		v1, err := vm.stack.Get(1)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v1); err != nil {
			return vm.runtimeError(err)
		}
	case OpPick:
		n := int(vm.code[vm.pc])
		vm.pc++
		v, err := vm.stack.Get(n)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpRoll:
		n := int(vm.code[vm.pc])
		vm.pc++
		v, err := vm.stack.PopAt(n)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpTuck:
		n := int(vm.code[vm.pc])
		vm.pc++
		v, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err = vm.stack.InsertAt(n, v); err != nil {
			return vm.runtimeError(err)
		}
	case OpRet:
		vm.runstate = RsComplete
	case OpFail:
		vm.runstate = RsError
		return vm.runtimeError(newRuntimeError("fail opcode invoked"))
	case OpZero:
		if err := vm.stack.Push(NewNumber(0)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPush1, OpPush2, OpPush3, OpPush4, OpPush5, OpPush6, OpPush7, OpPush8:
		// use a mask to retrieve the actual count of bytes to fetch
		nbytes := byte(instr) & 0xF
		var value int64
		var i byte
		var b Opcode
		for i = 0; i < nbytes; i++ {
			b = vm.code[vm.pc]
			vm.pc++
			value |= int64(b) << (i * 8)
		}
		if b&0x80 != 0 {
			for i := nbytes; i < 8; i++ {
				value |= 0xFF
			}
		}
		if err := vm.stack.Push(NewNumber(value)); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushT:
		var value int64
		var i byte
		var b Opcode
		for i = 0; i < 8; i++ {
			b = vm.code[vm.pc]
			vm.pc++
			value |= int64(b) << (i * 8)
		}
		ts := NewTimestamp(value)
		if err := vm.stack.Push(ts); err != nil {
			return vm.runtimeError(err)
		}
	case OpNow:
		ts, err := vm.now.Now()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(ts); err != nil {
			return vm.runtimeError(err)
		}
	case OpPushB, OpPushA:
		n := int(vm.code[vm.pc])
		vm.pc++
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = byte(vm.code[vm.pc])
			vm.pc++
		}
		v := NewBytes(b)
		if err := vm.stack.Push(v); err != nil {
			return vm.runtimeError(err)
		}
	case OpOne:
		if err := vm.stack.Push(NewNumber(1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNeg1:
		if err := vm.stack.Push(NewNumber(-1)); err != nil {
			return vm.runtimeError(err)
		}
	case OpRand:
		r, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(r)); err != nil {
			return vm.runtimeError(err)
		}

	case OpPushL:
		if err := vm.stack.Push(NewList()); err != nil {
			return vm.runtimeError(err)
		}
	case OpAdd, OpMul, OpDiv, OpMod:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n2, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch instr {
		case OpAdd:
			t, err = signed.Add(n2, n1)
			if err != nil {
				return vm.runtimeError(err)
			}
		case OpMul:
			t, err = signed.Mul(n2, n1)
			if err != nil {
				return vm.runtimeError(err)
			}
		case OpDiv:
			t, err = signed.Div(n2, n1)
			if err != nil {
				return vm.runtimeError(err)
			}
		case OpMod:
			t, err = signed.Mod(n2, n1)
			if err != nil {
				return vm.runtimeError(err)
			}
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpDivMod:
		d, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		q, m, err := signed.DivMod(n, d)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(m)); err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(q)); err != nil {
			return vm.runtimeError(err)
		}
	case OpMulDiv:
		d, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		v, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		v2, err := signed.MulDiv(v, n, d)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(v2)); err != nil {
			return vm.runtimeError(err)
		}
	case OpSub:
		// Subtraction is special because you can also subtract timestamps
		v1, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		v2, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch n1 := v1.(type) {
		case Number:
			n2, ok := v2.(Number)
			if !ok {
				return vm.runtimeError(newRuntimeError("incompatible types"))
			}
			t, err = signed.Sub(n2.v, n1.v)
		case Timestamp:
			n2, ok := v2.(Timestamp)
			if !ok {
				return vm.runtimeError(newRuntimeError("incompatible types"))
			}
			t, err = signed.Sub(int64(n2.t), int64(n1.t))
		}
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpNot, OpNeg, OpInc, OpDec:
		n1, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		var t int64
		switch instr {
		case OpNot:
			if n1 == 0 {
				t = 1
			}
		case OpNeg:
			t = -n1
		case OpInc:
			t = n1 + 1
		case OpDec:
			t = n1 - 1
		}
		if err := vm.stack.Push(NewNumber(t)); err != nil {
			return vm.runtimeError(err)
		}
	case OpEq, OpGt, OpLt:
		v2, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		v1, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		r, err := v1.Compare(v2)
		if err != nil {
			return vm.runtimeError(err)
		}
		var result int64
		switch instr {
		case OpEq:
			if r == 0 {
				result = 1
			}
		case OpGt:
			if r > 0 {
				result = 1
			}
		case OpLt:
			if r < 0 {
				result = 1
			}
		}
		if err := vm.stack.Push(NewNumber(result)); err != nil {
			return vm.runtimeError(err)
		}

	case OpIndex:
		n, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if n >= l.Len() || n < 0 {
			return vm.runtimeError(newRuntimeError("list index out of bounds"))
		}
		if err := vm.stack.Push(l[n]); err != nil {
			return vm.runtimeError(err)
		}

	case OpLen:
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(NewNumber(int64(len(l)))); err != nil {
			return vm.runtimeError(err)
		}

	case OpAppend:
		v, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		// Limit list size
		if l.Len()+1 > MaxListSize {
			return vm.runtimeError(newRuntimeError("resulting list too large"))
		}
		newlist := l.Append(v)
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpExtend:
		l1, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		l2, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		// Limit list size
		if l1.Len()+l2.Len() > MaxListSize {
			return vm.runtimeError(newRuntimeError("resulting list too large"))
		}

		newlist := l2.Extend(l1)
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpSlice:
		end, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		begin, err := vm.stack.PopAsInt64()
		if err != nil {
			return vm.runtimeError(err)
		}
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if begin < 0 || begin > l.Len() || end < 0 || end > l.Len() || begin > end {
			return vm.runtimeError(newRuntimeError("index out of range in slice"))
		}
		newlist := l[begin:end]
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpField:
		st, err := vm.stack.PopAsStruct()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		f, err := st.Field(int(fix))
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(f); err != nil {
			return vm.runtimeError(err)
		}

	case OpFieldL:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		extract := func(v Value) (Value, error) {
			f, ok := v.(Struct)
			if !ok {
				return v, errors.New("fieldl requires list of non-struct")
			}
			return f.Field(int(fix))
		}
		result, err := src.Map(extract)
		if err != nil {
			return vm.runtimeError(err)
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpDef:
		// if we try to execute a def statement there has been an error and we should abort
		return vm.runtimeError(newRuntimeError("tried to execute def opcode"))

	case OpCall:
		// The call opcode tracks the number of the current routine, and will only call a
		// function that has a larger number than itself (this prevents recursion). It constructs a
		// new stack for the function and populates it by copying (NOT popping off!) the specified number of
		// values from the caller's stack. The function call returns a single Value which is pushed
		// onto the caller's stack.
		funcnum := int(vm.code[vm.pc])
		nargs := int(vm.code[vm.pc+1])
		vm.pc += 2
		result, err := vm.callFunction(funcnum, nargs, debug)
		if err != nil {
			return err
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpDeco:
		funcnum := int(vm.code[vm.pc])
		nargs := int(vm.code[vm.pc+1])
		vm.pc += 2
		// we're going to iterate over a List of structs so validate it
		l, err := vm.stack.PopAsListOfStructs(-1)
		if err != nil {
			return vm.runtimeError(err)
		}
		newlist := NewList()
		for i := range l {
			// This is safe because we checked above
			s, _ := l[i].(Struct)
			retval, err := vm.callFunction(funcnum, nargs, debug, s)
			if err != nil {
				return err
			}
			// in order to prevent memory bombs, deco cannot add non-scalars
			if !retval.IsScalar() {
				return vm.runtimeError(newRuntimeError("deco result must be scalar"))
			}
			newlist = newlist.Append(s.Append(retval))
		}
		if err := vm.stack.Push(newlist); err != nil {
			return vm.runtimeError(err)
		}

	case OpEndDef:
		// we hit this at the end of a function that hasn't used OpRet or OpFail
		vm.runstate = RsComplete

	case OpIfZ:
		t, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		success := false
		switch v := t.(type) {
		case Number:
			success = (v.v == 0)
		default:
			// nothing else can test true for zeroness
		}
		// if we did not succeed, we have to skip to the matching end or else opcode
		if !success {
			if err := vm.skipToMatchingBracket(); err != nil {
				return vm.runtimeError(err)
			}
		}
	case OpIfNZ:
		t, err := vm.stack.Pop()
		if err != nil {
			return vm.runtimeError(err)
		}
		success := false
		switch v := t.(type) {
		case Number:
			success = (v.v != 0)
		default:
			// everything else tests true for nonzeroness
			success = true
		}
		// if we did not succeed, we have to skip to the matching END or ELSE opcode
		if !success {
			if err := vm.skipToMatchingBracket(); err != nil {
				return vm.runtimeError(err)
			}
		}
	case OpElse:
		// if we hit this in execution, it means we did the first clause of an if statement
		// and now need to skip to the matching end
		if err := vm.skipToMatchingBracket(); err != nil {
			return vm.runtimeError(err)
		}
	case OpEndIf:
		// OpEndIf is a no-op (it is only hit when it ends an if block that evaluated to true
		// and there was no Else clause

	case OpSum, OpAvg, OpMax, OpMin:
		l, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}

		// Define helper functions for the Reduce function
		sum := func(prev, current Value) Value {
			// prev is guaranteed to be a Number, but current might not be
			c, ok := current.(Number)
			if !ok {
				return prev
			}
			return NewNumber(prev.(Number).v + c.v)
		}
		max := func(prev, current Value) Value {
			cmp, _ := prev.Compare(current)
			if cmp < 0 {
				return current
			}
			return prev
		}
		min := func(prev, current Value) Value {
			cmp, _ := current.Compare(prev)
			if cmp < 0 {
				return current
			}
			return prev
		}

		var result Value
		switch instr {
		case OpSum:
			result = l.Reduce(sum, NewNumber(0))
		case OpAvg:
			if l.Len() == 0 {
				return vm.runtimeError(newRuntimeError("cannot average empty list"))
			}
			result = NewNumber(l.Reduce(sum, NewNumber(0)).(Number).v / l.Len())
		case OpMin:
			result = l.Reduce(min, NewNumber(math.MaxInt64))
		case OpMax:
			result = l.Reduce(max, NewNumber(math.MinInt64))
		}
		if err := vm.stack.Push(result); err != nil {
			return vm.runtimeError(err)
		}

	case OpChoice:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		if src.Len() == 0 {
			return vm.runtimeError(newRuntimeError("cannot use choice on empty list"))
		}
		i, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}
		r := rand.New(rand.NewSource(i))
		n := r.Intn(int(src.Len()))
		item := src[n]
		if err := vm.stack.Push(item); err != nil {
			return vm.runtimeError(err)
		}

	case OpWChoice:
		fix := vm.code[vm.pc]
		vm.pc++
		src, err := vm.stack.PopAsListOfStructs(int(fix))
		if err != nil {
			return vm.runtimeError(err)
		}

		if src.Len() == 0 {
			return vm.runtimeError(newRuntimeError("cannot use wchoice on an empty list"))
		}
		// because PopAsListOfStructs() validates the data,
		// we know we're safe to traverse the list of structs
		// and pull out our specified field as a Number
		sum := func(prev, current Value) Value {
			p := prev.(Number).v
			fi, _ := current.(Struct).Field(int(fix))
			c := fi.(Number).v
			return NewNumber(p + c)
		}
		total := src.Reduce(sum, NewNumber(0)).(Number).AsInt64()

		rand, err := vm.rand.RandInt()
		if err != nil {
			return vm.runtimeError(err)
		}

		var partialSum int64
		for i := range src {
			fi, _ := src[i].(Struct).Field(int(fix))
			partialSum += fi.(Number).AsInt64()
			if FractionLess(rand, math.MaxInt64, partialSum, total) {
				err := vm.stack.Push(src[i])
				if err != nil {
					return vm.runtimeError(err)
				}
				return nil
			}
		}

		// if we get here, something is very wrong
		panic(fmt.Sprintf("wchoice can't happen: %d %d %d", rand, partialSum, total))

	case OpSort:
		src, err := vm.stack.PopAsList()
		if err != nil {
			return vm.runtimeError(err)
		}
		fix := vm.code[vm.pc]
		vm.pc++
		// note - error handling is weak because the less function that sort.Slice()
		// uses cannot fail, so we can only figure it out after the sort completes.
		// This means if you try to sort bad data, you still get an error but
		// the sort finishes first.
		hadErr := false
		less := func(i, j int) bool {
			si, ok1 := src[i].(Struct)
			sj, ok2 := src[j].(Struct)
			if !ok1 || !ok2 {
				hadErr = true
				return false
			}
			fi, e1 := si.Field(int(fix))
			fj, e2 := sj.Field(int(fix))
			cmp, e3 := fi.Compare(fj)
			if e1 != nil || e2 != nil || e3 != nil {
				hadErr = true
			}
			return cmp == -1
		}
		sort.Slice(src, less)
		if hadErr {
			return vm.runtimeError(newRuntimeError("sort error"))
		}
		if err := vm.stack.Push(src); err != nil {
			return vm.runtimeError(err)
		}

	case OpLookup:
		funcnum := int(vm.code[vm.pc])
		nargs := int(vm.code[vm.pc+1])
		vm.pc += 2
		// we're going to iterate over a List of structs so validate it
		l, err := vm.stack.PopAsListOfStructs(-1)
		if err != nil {
			return vm.runtimeError(err)
		}
		foundix := -1
		for i := range l {
			// This is safe because we checked above
			s, _ := l[i].(Struct)
			result, err := vm.callFunction(funcnum, nargs, debug, s)
			if err != nil {
				return err
			}
			if n, ok := result.(Number); ok {
				if n.AsInt64() != 0 {
					foundix = i
					break
				}
			}
		}
		if foundix == -1 {
			return vm.runtimeError(errors.New("lookup failed"))
		}
		if err := vm.stack.Push(NewNumber(int64(foundix))); err != nil {
			return vm.runtimeError(err)
		}

	default:
		return vm.runtimeError(newRuntimeError(fmt.Sprintf("unimplemented opcode %s at %d", instr, vm.pc)))
	}

	return nil
}
