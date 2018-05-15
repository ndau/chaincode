package main

import (
	"bytes"
	"io"
	"log"
	"os"

	arg "github.com/alexflint/go-arg"
	"github.com/oneiro-ndev/chaincode/pkg/vm"
)

func main() {
	var args struct {
		Input   string `arg:"positional"`
		Output  string `arg:"-o" help:"Output filename"`
		Comment string `arg:"-c" help:"Comment to embed in the output file."`
	}
	arg.MustParse(&args)

	name := "stdin"
	in := os.Stdin
	if args.Input != "" {
		name = args.Input
		f, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	var buf bytes.Buffer
	tee := io.TeeReader(in, &buf)

	sn, err := ParseReader("", tee)
	if err != nil {
		log.Fatal(describeErrors(err, buf.String()))
	}

	out := os.Stdout
	if args.Output != "" {
		f, err := os.Create(args.Output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out = f
	}

	b := sn.(*Script).bytes()
	err = vm.Serialize(name, args.Comment, b, out)
	if err != nil {
		log.Fatal(err)
	}
}
