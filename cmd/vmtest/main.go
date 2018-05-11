package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

func main() {
	var args struct {
		Input   string   `arg:"-i" help:"File to read from (default stdin)"`
		History bool     `arg:"-h" help:"Dump history after running"`
		Stack   []string `arg:"positional" help:"Values to put on the stack, topmost last"`
	}
	arg.MustParse(&args)

	in := os.Stdin
	if args.Input != "" {
		fname := args.Input
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	binary, err := vm.Deserialize(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Name: %s\n     '%s'\n", binary.Name, binary.Comment)
	machine, err := vm.New(binary)
	if err != nil {
		log.Fatal(err)
	}
	startingStack := []vm.Value{}
	for _, s := range args.Stack {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		n := vm.NewNumber(v)
		startingStack = append(startingStack, n)
	}
	machine.Init(startingStack)
	fmt.Println(machine.Stack())
	fmt.Println("Running")
	err = machine.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Complete")
	fmt.Println(machine.Stack())
	if args.History {
		fmt.Println("-- History --")
		for _, h := range machine.History() {
			st := strings.Split(h.Stack.String(), "\n")
			st1 := make([]string, len(st))
			for i := range st {
				st1[i] = st[i][4:]
			}
			fmt.Printf("PC: %3d STK: %s\n", h.PC, strings.Join(st1, ", "))
		}
	}
}
