package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oneiro-ndev/chaincode/pkg/chain"
	"github.com/oneiro-ndev/chaincode/pkg/vm"
	"github.com/oneiro-ndev/ndau/pkg/ndau/backing"
	"github.com/oneiro-ndev/ndaumath/pkg/types"

	arg "github.com/alexflint/go-arg"
)

// crank is a repl for chaincode

// crank -i inputstream -f FILE.chbin

// crank starts up and creates a new VM with no contents
// If -f was specified, crank attempts to load the given file instead of starting with an empty vm
// if -i was specified, crank then attempts to read the input file as if it were a series of
// commands. When -i terminates, it returns control to the normal input. If you want crank to
// terminate automatically, make sure your input file ends in a quit command.

// command is a type that is used to create a table of commands for the repl
type command struct {
	parms   string
	aliases []string
	summary string
	detail  string
	handler func(rs *runtimeState, args string) error
}

func (c command) matchesAlias(s string) bool {
	for _, a := range c.aliases {
		if s == a {
			return true
		}
	}
	return false
}

var commands = map[string]command{
	"help": command{
		aliases: []string{"?"},
		summary: "prints this help message (help verbose for extended explanation)",
		detail:  ``,
		handler: nil, //  we need to fill this in dynamically because it traverses this list
	},
	"quit": command{
		aliases: []string{"q"},
		summary: "exits the chain program",
		detail:  `Ctrl-D also works`,
		handler: func(rs *runtimeState, args string) error {
			fmt.Println("Goodbye.")
			os.Exit(0)
			return nil
		},
	},
	"load": command{
		aliases: []string{"l"},
		summary: "loads the file FILE as a chasm binary (.chbin)",
		detail:  `File must conform to the chasm binary standard.`,
		handler: (*runtimeState).load,
	},
	"run": command{
		aliases: []string{"r"},
		summary: "runs the currently loaded VM from the current IP",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			return rs.run(false)
		},
	},
	"next": command{
		aliases: []string{"n"},
		summary: "executes one opcode at the current IP and prints the status",
		detail:  `If the opcode is a function call, this executes the entire function call before stopping.`,
		handler: func(rs *runtimeState, args string) error {
			return rs.step(true)
		},
	},
	"trace": command{
		aliases: []string{"tr", "t"},
		summary: "runs the currently loaded VM from the current IP",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			return rs.run(true)
		},
	},
	"event": command{
		aliases: []string{"ev", "e"},
		summary: "sets the ID of the event to be executed (may change the current IP)",
		detail:  ``,
		handler: (*runtimeState).setevent,
	},
	"disassemble": command{
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
	"reset": command{
		aliases: []string{"k"},
		summary: "resets the VM to the event and stack that were current at the last Run, Trace, Push, Pop, or Event command",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			rs.reinit(rs.stack)
			fmt.Println(rs.vm.Stack())
			return nil
		},
	},
	"stack": command{
		aliases: []string{"k"},
		summary: "prints the contents of the stack",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			fmt.Println(rs.vm.Stack())
			return nil
		},
	},
	"push": command{
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
			for _, v := range topush {
				rs.vm.Stack().Push(v)
			}
			fmt.Println(rs.vm.Stack())
			return rs.reinit(rs.vm.Stack())
		},
	},
	"pop": command{
		aliases: []string{"o"},
		summary: "pops the top stack item and prints it",
		detail:  ``,
		handler: func(rs *runtimeState, args string) error {
			v, err := rs.vm.Stack().Pop()
			if err != nil {
				return err
			}
			fmt.Println(v)
			return rs.reinit(rs.vm.Stack())
		},
	},
}

type runtimeState struct {
	vm    *vm.ChaincodeVM
	event byte
	stack *vm.Stack
}

// getRandomAccount randomly generates an account object
// it probably needs to be smarter than this
func getRandomAccount() backing.AccountData {
	const ticksPerDay = 24 * 60 * 60 * 1000000
	t, _ := types.TimestampFrom(time.Now())
	ad := backing.NewAccountData(t)
	// give it a balance between .1 and 100 ndau
	ad.Balance = types.Ndau((rand.Intn(1000) + 1) * 1000000)
	// set WAA to some time within 45 days
	ad.WeightedAverageAge = types.Duration(rand.Intn(ticksPerDay * 45))

	ad.LastEAIUpdate = t.Add(types.Duration(-rand.Intn(ticksPerDay * 3)))
	ad.LastWAAUpdate = t.Add(types.Duration(-rand.Intn(ticksPerDay * 10)))
	return ad
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
	// fields for structs are fieldid:Value; they are returned as a struct with one field that
	// is consolidated when they are enclosed in {} wrappers
	strfieldp := regexp.MustCompile("^([0-9]+) *:")

	s = strings.TrimSpace(s)
	retval := make([]vm.Value, 0)
	for s != "" {
		switch {
		case tsp.FindString(strings.ToUpper(s)) != "":
			t, err := vm.ParseTimestamp(tsp.FindString(strings.ToUpper(s)))
			if err != nil {
				panic(err)
			}
			s = s[len(tsp.FindString(strings.ToUpper(s))):]
			retval = append(retval, t)

		case strings.HasPrefix(s, "account"):
			str, err := chain.ToValue(getRandomAccount())
			s = s[len("account"):]
			if err != nil {
				return retval, err
			}
			retval = append(retval, str)

		case strings.HasPrefix(s, "["):
			if !strings.HasSuffix(s, "]") {
				return retval, errors.New("list start with no list end")
			}
			contents, err := parseValues(s[1 : len(s)-1])
			if err != nil {
				return retval, err
			}
			retval = append(retval, vm.NewList(contents...))
			// there can be only one list per line and it must end the line
			return retval, nil

		case strings.HasPrefix(s, "{"):
			if !strings.HasSuffix(s, "}") {
				return nil, errors.New("struct start with no struct end")
			}
			contents, err := parseValues(s[1 : len(s)-1])
			if err != nil {
				return nil, err
			}
			str := vm.NewStruct()
			for _, v := range contents {
				vs, ok := v.(*vm.Struct)
				if !ok {
					return retval, errors.New("untagged field in struct definition")
				}
				for _, ix := range vs.Indices() {
					v2, _ := vs.Get(ix)
					str = str.Set(ix, v2)
				}
			}
			return []vm.Value{str}, nil

		case strfieldp.FindString(s) != "":
			subm := strfieldp.FindStringSubmatch(s)
			f := subm[1]
			fieldid, _ := strconv.ParseInt(f, 10, 8)
			s = s[len(subm[0]):]
			contents, err := parseValues(s)
			if err != nil {
				return retval, err
			}
			if len(contents) == 0 {
				return retval, errors.New("field index without field value")
			}
			str := vm.NewStruct().Set(byte(fieldid), contents[0])
			retval = append(append(retval, str), contents[1:]...)
			return retval, nil

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

		default:
			return nil, errors.New("unparseable " + s)
		}
		s = strings.TrimSpace(s)
	}
	return retval, nil
}

func help(rs *runtimeState, args string) error {
	keys := make([]string, 0, len(commands))
	for key := range commands {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))
	for _, key := range keys {
		fmt.Printf("%11s: %s %s\n", key, commands[key].summary, commands[key].aliases)
		if strings.HasPrefix(args, "v") {
			fmt.Println(commands[key].detail)
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
func (rs *runtimeState) reinit(stk *vm.Stack) error {
	// copy the current stack and save it in case we need to reset
	rs.stack = stk.Clone()

	// now initialize
	return rs.vm.InitFromStack(rs.event, rs.stack)
}

// setevent sets up the VM to run the given event, which means that it calls
// reinit to set up the stack as well.
func (rs *runtimeState) setevent(eventid string) error {
	ev, err := strconv.ParseInt(strings.TrimSpace(eventid), 0, 8)
	if err != nil {
		return err
	}
	rs.event = byte(ev)

	return rs.reinit(rs.vm.Stack())
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

func (rs *runtimeState) dispatch(s string) error {
	p := regexp.MustCompile("[[:space:]]+")
	args := p.Split(s, 2)
	for key, cmd := range commands {
		if key == args[0] || cmd.matchesAlias(args[0]) {
			return cmd.handler(rs, args[1])
		}
	}
	return fmt.Errorf("unknown command %s - type ? for help", s)
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
	h := commands["help"]
	h.handler = help
	commands["help"] = h

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
