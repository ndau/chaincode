package vm

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"
)

type option struct {
	w int
	v interface{}
}

type weightings struct {
	opts  []option
	total int
}

// choose does a weighted selection of a single value from a collection of weighted values
func choose(w weightings) interface{} {
	if w.total == 0 {
		for i := 0; i < len(w.opts); i++ {
			w.total += w.opts[i].w
		}
	}
	r := rand.Intn(w.total)
	t := 0
	for i := 0; i < len(w.opts); i++ {
		t += w.opts[i].w
		if t >= r {
			return w.opts[i].v
		}
	}
	// we should never get here
	return w.opts[len(w.opts)-1].v
}

func randByte() byte {
	n := rand.Intn(256)
	return byte(n)
}

// genRandomProgram generates a program that will pass the VM's validation
// criteria, so we can make sure that the runtime doesn't just die when presented
// with any strange code. It generates a set of functions, and then
// generates a set of handlers that may call those functions.
func genRandomProgram() []string {
	// here are some weightings for the number of handlers in a program
	whandlers := weightings{opts: []option{
		option{50, 1},
		option{30, 2},
		option{10, 3},
		option{10, 4},
	}}
	// here are some weightings for the number of functions in a program
	wfuncts := weightings{opts: []option{
		option{50, 1},
		option{30, 2},
		option{10, 3},
		option{10, 4},
	}}

	nfuncs := choose(wfuncts).(int)
	s := []string{}
	for i := 0; i < nfuncs; i++ {
		s = append(s, genRandomFunc(i, nfuncs-1)...)
	}

	nhandlers := choose(whandlers).(int)
	for i := 0; i < nhandlers; i++ {
		handlerID := rand.Intn(256) // this will sometimes generate duplicated handler IDs, that's OK
		s = append(s, genRandomHandler(handlerID, nfuncs-1)...)
	}
	return s
}

// genRandomHandler generates a handler with the handler number hnum
// and also accepts a parameter for the maximum function number it can call.
func genRandomHandler(hnum, fmax int) []string {
	// here are some weightings for the number of top-level sequences in a handler
	w := weightings{opts: []option{
		option{40, 1},
		option{20, 2},
		option{10, 3},
		option{10, 4},
		option{5, 5},
		option{5, 6},
	}}
	s := []string{fmt.Sprintf("%s 1 %02x", OpHandler, hnum)}
	nseqs := choose(w).(int)
	for i := 0; i < nseqs; i++ {
		s = append(s, genRandomSequence(0, fmax, 0)...)
	}
	// check for empty results  -- that's useless to us
	// so if we did that, just try again
	if len(s) == 1 {
		return genRandomHandler(hnum, fmax)
	}
	s = append(s, OpEndDef.String())
	return s
}

// genRandomFunc generates a function with the function number fnum
// and also accepts a parameter for the maximum function number it can call.
func genRandomFunc(fnum, fmax int) []string {
	// here are some weightings for the number of top-level sequences in a function
	seqw := weightings{opts: []option{
		option{40, 1},
		option{20, 2},
		option{10, 3},
		option{10, 4},
		option{5, 5},
		option{5, 6},
	}}
	argw := weightings{opts: []option{
		option{40, 0},
		option{30, 1},
		option{15, 2},
		option{10, 3},
		option{5, 4},
	}}
	numargs := choose(argw).(int)
	s := []string{fmt.Sprintf("%s %02x %02x", OpDef, fnum, numargs)}

	nseqs := choose(seqw).(int)
	for i := 0; i < nseqs; i++ {
		s = append(s, genRandomSequence(fnum+1, fmax, 0)...)
	}
	// check for empty functions at level 0  -- that's useless to us
	// so if we did that, just try again
	if fnum == 0 && len(s) == 1 {
		return genRandomFunc(fnum, fmax)
	}
	s = append(s, OpEndDef.String())
	return s
}

// genRandomSequence creates a continuous code sequence, which may include
// other sequences. It accepts a range of function numbers that it is allowed to
// call; if min > max then no function calls will be generated. It won't go more than 3 ifs deep.
//
func genRandomSequence(fmin, fmax, depth int) []string {
	for {
		// here are some weightings for the different kinds of sequences
		seqw := weightings{opts: []option{
			option{50, OpNop},
			option{50, OpIfNZ},
			option{50, OpIfZ},
			option{10, OpCall},
			option{10, OpDeco},
			option{10, OpLookup},
		}}
		s := []string{}
		op := choose(seqw).(Opcode)
		switch op {
		case OpNop:
			return []string{genLinearSequence()}
		case OpIfZ, OpIfNZ:
			if depth > 3 {
				continue
			}
			s = append(s, op.String())
			s = append(s, genRandomSequence(fmin, fmax, depth+1)...)
			// sometimes we want an else clause
			if rand.Intn(100) < 30 {
				s = append(s, OpElse.String())
				s = append(s, genRandomSequence(fmin, fmax, depth+1)...)
			}
			s = append(s, OpEndIf.String())
			return s
		case OpCall, OpDeco, OpLookup:
			// can't do these if we're down too far in the call stack
			if fmin > fmax {
				continue
			}
			var fnum int
			if fmax > fmin {
				fnum = rand.Intn(fmax-fmin) + fmin
			} else {
				fnum = fmin
			}
			s = append(s, fmt.Sprintf("%s %02x", op, fnum))
		}
	}
}

// genLinearSequence generates a sequence of simple opcodes that executes linearly
// It might be empty.
func genLinearSequence() string {
	s := []string{}
	for rand.Intn(10) != 0 {
		s = append(s, genUnorderedInstruction())
	}
	return strings.Join(s, " ")
}

// genUnorderedInstruction creates individual instructions that are
// plausible (right number of bytes, etc). It only generates one instruction
// at a time -- there is no attempt to make sure they have a plausible
// sequence.
//
// Opcodes that are part of multi-instruction sequences (if, def, call, etc)
// are excluded.
func genUnorderedInstruction() string {
	for {
		op := Opcode(randByte())
		if !EnabledOpcodes.Get(byte(op)) {
			continue
		}
		s := []string{op.String()}

		switch op {
		case OpDef, OpHandler, OpEndDef, OpCall, OpDeco, OpLookup, OpIfNZ, OpIfZ, OpElse, OpEndIf:
			continue
		case OpPushB, OpPushA:
			// append up to 16 extra bytes
			extra := rand.Intn(16)
			s = append(s, fmt.Sprintf("%02x", extra))
			for i := 0; i < extra; i++ {
				s = append(s, fmt.Sprintf("%02x", randByte()))
			}
		default:
			extra := extraBytes([]Opcode{op}, 0)
			for i := 0; i < extra; i++ {
				s = append(s, fmt.Sprintf("%02x", randByte()))
			}
		}
		return strings.Join(s, " ")
	}
}

type valueType int

const (
	vtUnset         valueType = iota
	vtNumber        valueType = iota
	vtBytes         valueType = iota
	vtList          valueType = iota
	vtTimestamp     valueType = iota
	vtStruct        valueType = iota
	vtListOfStructs valueType = iota
)

// genRandomValue generates a single randomized Value object
// If it happens to create a list, all the elements in the list will
// be of the same format, including a list of structs
func genRandomValue(vt valueType) Value {
	alltypes := weightings{opts: []option{
		option{50, vtNumber},
		option{30, vtList},
		option{20, vtStruct},
		option{20, vtListOfStructs},
		option{10, vtBytes},
		option{10, vtTimestamp},
	}}

	scalars := weightings{opts: []option{
		option{50, vtNumber},
		option{10, vtBytes},
		option{10, vtTimestamp},
	}}

	if vt == vtUnset {
		vt = choose(alltypes).(valueType)
	}
	switch vt {
	case vtNumber:
		return NewNumber(rand.Int63())
	case vtBytes:
		nbytes := rand.Intn(16)
		b := make([]byte, nbytes)
		for i := 0; i < nbytes; i++ {
			// b[i] = randByte()
			b[i] = byte(rand.Intn(26) + 97)
		}
		return NewBytes(b)
	case vtList:
		nitems := rand.Intn(16)
		itemtype := choose(scalars).(valueType)
		l := make([]Value, nitems)
		for i := 0; i < nitems; i++ {
			l[i] = genRandomValue(itemtype)
		}
		return NewList(l...)
	case vtTimestamp:
		return NewTimestampFromInt(rand.Int63())
	case vtStruct:
		// create a struct of 1-5 scalar members
		vs := genRandomValues(1, 5)
		return NewTestStruct(vs...)
	case vtListOfStructs:
		// number of items in the list
		nitems := rand.Intn(16)
		// number of fields per struct
		nfields := rand.Intn(4) + 1

		// initialize the struct
		l := make([]Value, nitems)
		for i := 0; i < nitems; i++ {
			l[i] = NewStruct()
		}

		// now populate it
		for f := 0; f < nfields; f++ {
			ft := choose(scalars).(valueType)
			for i := 0; i < nitems; i++ {
				l[i] = l[i].(*Struct).Set(byte(f), genRandomValue(ft))
			}
		}
		return NewList(l...)
	}
	// we should never get here
	return NewNumber(0)
}

// genRandomValues creates an array of between min and max random Value objects (inclusive);
// if zeroOk is true, it might not create any at all.
func genRandomValues(min, max int) []Value {
	var nValues = rand.Intn(max+1-min) + min
	ret := []Value{}
	for i := 0; i < nValues; i++ {
		ret = append(ret, genRandomValue(vtUnset))
	}
	return ret
}

// key reads an error that may end with extra data and returns only
// the leading textual part of it so that we can aggregate messages
// that are similar but not identical. I wouldn't do it this way
// in production but for testing it's fine.
func key(err error) string {
	s := err.Error()
	p := regexp.MustCompile("^[a-zA-Z '-]+")
	return p.FindString(s)
}

// TestFuzzHandlers generates individual handlers of random enabled opcodes -- they have
// the proper wrappers at the beginning and end but are otherwise random-but-valid opcodes.
// Many of these will just fail validation.
func TestFuzzHandlers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var prog string
	var savedvm *ChaincodeVM
	var startingStack []Value
	// we want to know what failed if something failed
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Run caused a panic:", r)
			fmt.Println("Program: ", prog)
			fmt.Println("Starting stack:", startingStack)
			fmt.Println(savedvm)
			debug.PrintStack()
			t.Errorf("Test failed.")
		}
	}()

	rand.Seed(time.Now().UnixNano())
	results := make(map[string]int)
	total := 100
	nruns := os.Getenv("FUZZ_RUNS")
	if nruns != "" {
		total, _ = strconv.Atoi(nruns)
	}
	for i := 0; i < total; i++ {
		s := []string{OpHandler.String(), " 00"}
		for j := 0; j < rand.Intn(20)+5; j++ {
			op := Opcode(randByte())
			if !EnabledOpcodes.Get(byte(op)) {
				continue
			}
			s = append(s, op.String())
		}
		s = append(s, OpEndDef.String())
		prog = strings.Join(s, " ")
		ops := MiniAsm(prog)
		bin := ChasmBinary{"test", "TEST", ops}
		vm, err := New(bin)
		if err != nil {
			// fmt.Printf("Didn't load because %s: %s\n", err, p)
			results[key(err)]++
			continue
		}
		savedvm = vm

		// Put a couple of items on the stack
		startingStack = genRandomValues(1, 3)
		vm.Init(0, startingStack...)
		err = vm.Run(false)
		if err == nil {
			// fmt.Printf("Successfully ran:\n")
			// vm.DisassembleAll()
			results["success"]++
		} else {
			results[key(err)]++
		}
	}
	fmt.Printf("Results for %d runs:\n", total)
	for k, v := range results {
		fmt.Printf("%8d: %s\n", v, k)
	}
}

// TestFuzzValid generates programs that will pass validation (provided they're not too long).
// Because validation is so picky, this lets us have a much higher hit rate of programs that will
// actually exercise opcodes.
func TestFuzzValid(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var prog string
	var savedvm *ChaincodeVM
	var startingStack []Value
	// we want to know what failed if something failed
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Run caused a panic:", r)
			fmt.Println("Program: ", prog)
			fmt.Println("Starting stack:", startingStack)
			fmt.Println(savedvm)
			debug.PrintStack()
			t.Errorf("Test failed.")
		}
	}()

	rand.Seed(time.Now().UnixNano())
	results := make(map[string]int)
	total := 1000
	nruns := os.Getenv("FUZZ_RUNS")
	if nruns != "" {
		total, _ = strconv.Atoi(nruns)
	}
	fmt.Printf("Running %d iterations.\n", total)
	attempts := 0
	secondaryFailures := 0
	for attempts < total {
		prog = strings.Join(genRandomProgram(), "\n")
		ops := MiniAsm(prog)
		bin := ChasmBinary{"test", "TEST", ops}
		vm, err := New(bin)
		if err != nil {
			// fmt.Printf("Didn't load because %s: %s\n", err, vm)
			results[key(err)]++
			continue
		}
		savedvm = vm

		// try the script at least once for each handler
		events := vm.HandlerIDs()
		var h byte
		for hix := 0; hix < len(events); hix++ {
			h = byte(events[hix])
		}
		// if the script runs to completion, try it a few more times with some other
		// data to see if we can make it fail
		for runcount := 0; runcount < 10; runcount++ {
			attempts++
			// Put a couple of items on the stack
			startingStack = genRandomValues(1, 3)
			vm.Init(h, startingStack...)
			err = vm.Run(false)
			if err != nil {
				results[key(err)]++
				if runcount > 1 {
					secondaryFailures++
				}
				break
			}
			// fmt.Printf("Successfully ran:\n")
			// vm.DisassembleAll()
			results["success"]++
		}
	}
	fmt.Printf("Results for %d runs (%d secondary failures):\n", attempts, secondaryFailures)
	for k, v := range results {
		fmt.Printf("%8d: %s\n", v, k)
	}
}

func TestFuzzJunk(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var prog string
	var savedvm *ChaincodeVM
	var startingStack []Value
	// we want to know what failed if something failed
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Run caused a panic:", r)
			fmt.Println("Program: ", prog)
			fmt.Println("Starting stack:", startingStack)
			fmt.Println(savedvm)
			debug.PrintStack()
			t.Errorf("Test failed.")
		}
	}()

	rand.Seed(time.Now().UnixNano())
	results := make(map[string]int)
	total := 100
	nruns := os.Getenv("FUZZ_RUNS")
	if nruns != "" {
		total, _ = strconv.Atoi(nruns)
	}
	for i := 0; i < total; i++ {
		s := []string{}
		// there's a 10% chance of not having a function starter
		if rand.Intn(100) < 10 {
			s = append(s, OpDef.String(), " 00")
		}
		for j := 0; j < rand.Intn(50)+5; j++ {
			b := randByte()
			s = append(s, fmt.Sprintf("%02x", b))
		}
		s = append(s, OpEndDef.String())
		prog = strings.Join(s, " ")
		ops := MiniAsm(prog)
		bin := ChasmBinary{"test", "TEST", ops}
		vm, err := New(bin)
		if err != nil {
			// fmt.Printf("Didn't load because %s: %s\n", err, p)
			results[key(err)]++
			continue
		}
		savedvm = vm

		// Put a couple of items on the stack
		startingStack = genRandomValues(1, 3)
		vm.Init(0, startingStack...)
		err = vm.Run(false)
		if err == nil {
			// fmt.Printf("Successfully ran:\n")
			// vm.DisassembleAll()
			results["success"]++
		} else {
			results[key(err)]++
		}
	}
	fmt.Printf("Results for %d runs:\n", total)
	for k, v := range results {
		fmt.Printf("%8d: %s\n", v, k)
	}
}
