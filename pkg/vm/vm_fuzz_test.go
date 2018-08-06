package vm

import (
	"fmt"
	"math/rand"
	"runtime/debug"
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

func genRandomProgram() []string {
	// here are some weightings for the number of functions in a program
	w := weightings{opts: []option{
		option{50, 1},
		option{30, 2},
		option{10, 3},
		option{10, 4},
	}}
	nfuncs := choose(w).(int)
	s := []string{}
	for i := 0; i < nfuncs; i++ {
		s = append(s, genRandomFunc(i, nfuncs-1)...)
	}
	return s
}

// genRandomFunc generates a function with the function number fnum
// and also accepts a parameter for the maximum function number it can call.
func genRandomFunc(fnum, fmax int) []string {
	// here are some weightings for the number of top-level sequences in a function
	w := weightings{opts: []option{
		option{40, 1},
		option{20, 2},
		option{10, 3},
		option{10, 4},
		option{5, 5},
		option{5, 6},
	}}
	s := []string{fmt.Sprintf("%s %02x", OpDef, fnum)}
	nseqs := choose(w).(int)
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
		argw := weightings{opts: []option{
			option{50, 0},
			option{30, 1},
			option{10, 2},
			option{5, 3},
			option{5, 4},
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
			numargs := choose(argw).(int)
			s = append(s, fmt.Sprintf("%s %02x %02x", op, fnum, numargs))
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
		if !EnabledOpcodes.Get(int(op)) {
			continue
		}
		s := []string{op.String()}

		switch op {
		case OpDef, OpEndDef, OpCall, OpDeco, OpLookup, OpIfNZ, OpIfZ, OpElse, OpEndIf:
			continue
		case OpPushB, OpPushA:
			// append up to 64 extra bytes
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

func TestFuzz(t *testing.T) {
	var prog string
	var savedvm *ChaincodeVM
	// we want to know what failed if something failed
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Run caused a panic:", r)
			fmt.Println("Program: ", prog)
			fmt.Println(savedvm)
			debug.PrintStack()
		}
	}()

	rand.Seed(time.Now().UnixNano())
	noasm := 0
	ranok := 0
	total := 10000
	for i := 0; i < total; i++ {
		prog = strings.Join(genRandomProgram(), "\n")
		ops := miniAsm(prog)
		bin := ChasmBinary{"test", "", "TEST", ops}
		vm, err := New(bin)
		if err != nil {
			// fmt.Printf("Didn't load because %s: %s\n", err, p)
			noasm++
			continue
		}
		savedvm = vm

		// Put a couple of items on the stack
		vm.Init(NewNumber(1), NewNumber(2))
		err = vm.Run(false)
		if err == nil {
			// fmt.Printf("Successfully ran:\n")
			// vm.DisassembleAll()
			ranok++
		}
	}
	fmt.Printf("Attempts: %d\nAsm Failures: %d\nRan OK: %d\n", total, noasm, ranok)
}
