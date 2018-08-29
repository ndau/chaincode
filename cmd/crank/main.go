package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/oneiro-ndev/chaincode/pkg/vm"

	arg "github.com/alexflint/go-arg"
)

// crank is a repl for chaincode

// crank -i inputstream -b FILE.chbin

// crank starts up and creates a new VM with no contents
// If -f was specified, crank attempts to load the given file instead of starting with an empty vm
// if -i was specified, crank then attempts to read the input file as if it were a series of
// commands. When -i terminates, it returns control to the normal input. If you want crank to
// terminate automatically, make sure your input file ends in a quit command.

// Types:
//    Number: (int, int64, uint64)
//    Bytes:   (string, []byte)
//    Timestamp: (time.Time)
//    Struct: struct of above (but not list)
//    List: ([] any of the above)

// type S struct {
// 	N int `chain:"0,name"
// }

// commands:
//     Load FILE
//         loads and validates a file
//     stacK
//         prints the current stack
//     Push value
//         pushes a new value on top of the current stack
//     pOp
//         pops the top of the stack and prints it
//     Run
//         runs the VM from current IP
//     Trace
//         runs the VM from the current IP with tracing on
//     Event EVT
//         sets the VM to run event EVT
//     Ip offset
//         sets the VM to start running at IP
//     input filename
//         reads filename and treats each line as if it were typed in at the command prompt
//     Step
//         single-steps the VM
//     Asm offset
//         starts the mini-assembler at the given offset; blank line to terminate
//         (offset is current IP if not otherwise stated)
//     Disasm offset
//         disassembles from the given offset (default current IP)
// 	   Quit
//         exits crank
//
// Value syntax:
//     Number (decimal, hex)
//     Timestamp
//     Quoted string
//     B(hex pairs)
//     [ list of values ] (commas or whitespace)
//     { list of values }(commas or whitespace)

type runtimeState struct {
	vm    *vm.ChaincodeVM
	event byte
}

func parseValues(s string) ([]vm.Value, error) {
	// timestamp
	tsp := regexp.MustCompile("^[0-9-]+T[0-9:]+Z")
	// address is 48 chars starting with nd and not containing io10
	// addrp := regexp.MustCompile("^nd[2-9a-km-np-zA-KM-NP-Z]{46}")
	// number is a base-10 signed integer OR a hex value starting with 0x
	nump := regexp.MustCompile("^0x([0-9A-Fa-f]+)|^-?[0-9]+")
	// hexp := regexp.MustCompile("^0x([0-9A-Fa-f]+)")
	// quoted strings can be either single, double, or triple quotes of either kind; no escapes.
	quotep := regexp.MustCompile(`^'([^']*)'|^"([^"]*)"|^'''(.*)'''|^"""(.*)"""`)
	// arrays of bytes are B(hex) with individual bytes as hex strings with no 0x; embedded spaces are ignored
	bytep := regexp.MustCompile(`^B\((([0-9A-Fa-f][0-9A-Fa-f] *)+)\)`)

	s = strings.TrimSpace(s)
	retval := make([]vm.Value, 0)
	for s != "" {
		switch {
		case strings.HasPrefix(s, "["):
			if !strings.HasSuffix(s, "]") {
				return retval, errors.New("list start with no list end")
			}
			contents, err := parseValues(s[1 : len(s)-1])
			if err != nil {
				return retval, err
			}
			retval = append(retval, vm.NewList(contents...))

			// case strings.HasPrefix(s, "{"):
			// 	if !strings.HasSuffix(s, "}") {
			// 		return nil, errors.New("struct start with no struct end")
			// 	}
			// 	contents, err := parseValues(s[1:len(s)-1])
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	str := vm.NewStruct()
			// 	for _, v := range contents {
			// 		str
			// 	}
			// 	return str, nil

		// case hexp.FindString(s) != "":
		// 	subm := hexp.FindStringSubmatch(s)
		// 	contents := subm[1]
		// 	s = s[len(subm[0]):]
		// 	n, _ := strconv.ParseInt(contents, 16, 64)
		// 	retval = append(retval, vm.NewNumber(n))

		case nump.FindString(s) != "":
			found := nump.FindString(s)
			s = s[len(found):]
			n, _ := strconv.ParseInt(found, 0, 64)
			retval = append(retval, vm.NewNumber(n))

		case bytep.FindString(s) != "":
			ba := []byte{}
			// the stream of bytes is the first submatch here
			submatches := bytep.FindStringSubmatch(s)
			contents := submatches[1]
			s = s[len(submatches[0]):]
			pair := regexp.MustCompile("([0-9A-Fa-f][0-9A-Fa-f])")
			for _, it := range pair.FindAllString(contents, -1) {
				b, _ := strconv.ParseInt(strings.TrimSpace(it), 16, 8)
				ba = append(ba, byte(b))
			}
			retval = append(retval, vm.NewBytes(ba))

		case quotep.FindString(s) != "":
			subm := quotep.FindSubmatch([]byte(s))
			contents := subm[1]
			s = s[len(subm[0]):]
			retval = append(retval, vm.NewBytes(contents))

		case tsp.FindString(strings.ToUpper(s)) != "":
			s = s[len(tsp.FindString(strings.ToUpper(s))):]
			t, err := vm.ParseTimestamp(strings.ToUpper(s))
			if err != nil {
				panic(err)
			}
			retval = append(retval, t)

		default:
			return nil, errors.New("unparseable " + s)
		}
		s = strings.TrimSpace(s)
	}
	return retval, nil
}

// load is a command that loads a file into a VM (or errors trying)
func (rs *runtimeState) load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	bin, err := vm.Deserialize(f)
	if err != nil {
		return err
	}
	vm, err := vm.New(bin)
	if err != nil {
		return err
	}
	vm.Init(0)
	rs.vm = vm
	return nil
}

func (rs *runtimeState) run(debug bool) error {
	d := rs.vm.Stack().Depth()
	values := make([]vm.Value, d)
	for i := 0; i < d; i++ {
		v, _ := rs.vm.Stack().Pop()
		values[d-i-1] = v
	}
	err := rs.vm.Init(rs.event, values...)
	if err != nil {
		return err
	}
	err = rs.vm.Run(debug)
	return err
}

func (rs *runtimeState) setevent(eventid string) error {
	ev, err := strconv.ParseInt(eventid, 0, 8)
	if err != nil {
		return err
	}
	rs.event = byte(ev)
	return nil
}

func (rs *runtimeState) dispatch(cmd string) error {
	p := regexp.MustCompile("[[:space:]]+")
	args := p.Split(cmd, -1)
	switch args[0] {
	case "quit", "q":
		fmt.Println("Goodbye.")
		os.Exit(0)
	case "load", "l":
		return rs.load(args[1])
	case "run", "r":
		return rs.run(false)
	case "trace", "t":
		return rs.run(true)
	// case "step", "s":
	// 	return rs.step()
	case "event", "ev", "e":
		return rs.setevent(args[1])
	case "dis", "d", "disasm", "disassemble":
		if rs.vm == nil {
			return errors.New("no VM is loaded")
		}
		rs.vm.DisassembleAll()
		return nil
	case "stack", "k":
		fmt.Println(rs.vm.Stack())
	case "push", "pu":
		s := strings.Join(args[1:], " ")
		topush, err := parseValues(s)
		if err != nil {
			return err
		}
		fmt.Println(topush)
		for _, v := range topush {
			rs.vm.Stack().Push(v)
		}
		fmt.Println(rs.vm.Stack())
	case "pop":
		v, err := rs.vm.Stack().Pop()
		if err != nil {
			return err
		}
		fmt.Println(v)
	default:
		return errors.New("unknown command " + cmd)
	}
	return nil
}

func (rs *runtimeState) repl(cmdsrc io.Reader) {
	reader := bufio.NewReader(os.Stdin)
	usingStdin := true
	if cmdsrc != nil {
		reader = bufio.NewReader(cmdsrc)
		usingStdin = false
	}
	for {
		fmt.Print("> ")
		s, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF && usingStdin == true {
			// we're really done now, shut down normally
			s = "quit\n"
			err = nil
		}
		if err == io.EOF {
			reader = bufio.NewReader(os.Stdin)
			fmt.Println("*** Input now from stdin ***")
			usingStdin = true
		}
		err = rs.dispatch(s)
		if err != nil {
			fmt.Println("  -> Error: ", err)
		}
		if rs.vm == nil {
			fmt.Println("  [no VM is loaded]")
		} else {
			rs.vm.Disassemble(rs.vm.IP())
		}
	}
}

func main() {
	var args struct {
		Input  string `arg:"-i" help:"Input command file"`
		Binary string `arg:"-b" help:"Binary file to load"`
	}
	arg.MustParse(&args)
	var inf io.Reader
	if args.Input != "" {
		var err error
		inf, err = os.Open(args.Input)
		if err != nil {
			panic(err)
		}
	}

	rs := runtimeState{}
	if args.Binary != "" {
		err := rs.load(args.Binary)
		if err != nil {
			panic(err)
		}
	}

	rs.repl(inf)
}
