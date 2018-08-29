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

// crank -i inputstream -f FILE.chbin

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

type command struct {
	name    string
	parms   string
	aliases []string
	summary string
	detail  string
	handler func(rs *runtimeState, args string) error
}

func (c command) matches(s string) bool {
	for _, a := range c.aliases {
		if s == a {
			return true
		}
	}
	return s == c.name
}

var commands = []command{
	command{
		name:    "help",
		aliases: []string{"?"},
		summary: "prints this help message (help verbose for extended explanation)",
		detail:  ``,
		handler: nil, //  we need to fill this in dynamically because it traverses this list
	},
	command{
		name:    "quit",
		aliases: []string{"q"},
		summary: "exits the chain program",
		detail:  `Ctrl-D also works`,
		handler: func(rs *runtimeState, args string) error {
			fmt.Println("Goodbye.")
			os.Exit(0)
			return nil
		},
	},
	command{
		name:    "load",
		aliases: []string{"l"},
		summary: "loads the file FILE as a chasm binary (.chbin)",
		detail:  `File must conform to the chasm binary standard.`,
		handler: (*runtimeState).load,
	},
	command{
		name:    "run",
		aliases: []string{"r"},
		summary: "runs the currently loaded VM from the current IP",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			return rs.run(false)
		},
	},
	command{
		name:    "step",
		aliases: []string{"s"},
		summary: "executes one opcode at the current IP and prints the status",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			return rs.step(true)
		},
	},
	command{
		name:    "trace",
		aliases: []string{"tr", "t"},
		summary: "runs the currently loaded VM from the current IP",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			return rs.run(true)
		},
	},
	command{
		name:    "event",
		aliases: []string{"ev", "e"},
		summary: "sets the ID of the event to be executed (may change the current IP)",
		detail:  ``,
		handler: (*runtimeState).setevent,
	},
	command{
		name:    "disassemble",
		aliases: []string{"dis", "disasm", "d"},
		summary: "disassembles the loaded vm",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			if rs.vm == nil {
				return errors.New("no VM is loaded")
			}
			rs.vm.DisassembleAll()
			return nil
		},
	},
	command{
		name:    "stack",
		aliases: []string{"k"},
		summary: "prints the contents of the stack",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			fmt.Println(rs.vm.Stack())
			return nil
		},
	},
	command{
		name:    "push",
		aliases: []string{"pu", "p"},
		summary: "pushes one or more values onto the stack",
		detail: `
Value syntax:
    Number (decimal, hex)
    Timestamp
    Quoted string
    B(hex pairs)
    [ list of values ] (commas or whitespace, must all be one line)
    { list of values }(commas or whitespace, must all be one line)
		`,
		handler: func(rs *runtimeState, args string) error {
			topush, err := parseValues(args)
			if err != nil {
				return err
			}
			fmt.Println(topush)
			for _, v := range topush {
				rs.vm.Stack().Push(v)
			}
			fmt.Println(rs.vm.Stack())
			return rs.reinit()
		},
	},
	command{
		name:    "pop",
		aliases: []string{"o"},
		summary: "pops the top stack item and prints it",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			v, err := rs.vm.Stack().Pop()
			if err != nil {
				return err
			}
			fmt.Println(v)
			return rs.reinit()
		},
	},
}

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

func help(rs *runtimeState, args string) error {
	for _, cmd := range commands {
		fmt.Printf("%11s: %s %s\n", cmd.name, cmd.summary, cmd.aliases)
		if strings.HasPrefix(args, "v") {
			fmt.Println(cmd.detail)
		}
	}
	return nil
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

// reinit calls init again, duplicating the entries that are currently
// on the stack.
func (rs *runtimeState) reinit() error {
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
	return nil
}

// setevent sets up the VM to run the given event, which means that it calls
// reinit to set up the stack as well.
func (rs *runtimeState) setevent(eventid string) error {
	ev, err := strconv.ParseInt(strings.TrimSpace(eventid), 0, 8)
	if err != nil {
		return err
	}
	rs.event = byte(ev)

	return rs.reinit()
}

func (rs *runtimeState) run(debug bool) error {
	err := rs.vm.Run(debug)
	return err
}

func (rs *runtimeState) step(debug bool) error {
	err := rs.vm.Step(debug)
	fmt.Println(rs.vm)
	return err
}

func (rs *runtimeState) dispatch(cmd string) error {
	p := regexp.MustCompile("[[:space:]]+")
	args := p.Split(cmd, 2)
	for _, cmd := range commands {
		if cmd.matches(args[0]) {
			return cmd.handler(rs, args[1])
		}
	}
	return fmt.Errorf("unknown command %s - type ? for help", cmd)
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
	// this needs to be filled in dynamically because the help function traverses
	// the commands list.
	commands[0].handler = help
	var args struct {
		Input string `arg:"-i" help:"Input command file"`
		File  string `arg:"-f" help:"File to load as a chasm binary (*.chbin)."`
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
	if args.File != "" {
		err := rs.load(args.File)
		if err != nil {
			panic(err)
		}
	}

	rs.repl(inf)
}
