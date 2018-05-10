package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"

	arg "github.com/alexflint/go-arg"
)

type ChasmBinary struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Context string `json:"context"`
	Data    []byte `json:"data"`
}

var Contexts = map[byte]string{
	CtxTest:        "TEST",
	CtxNodePayout:  "NODE_PAYOUT",
	CtxEaiTiming:   "EAI_TIMING",
	CtxNodeQuality: "NODE_QUALITY",
	CtxMarketPrice: "MARKET_PRICE",
}

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
		name := args.Input
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

	b := sn.(Script).bytes()
	output := ChasmBinary{
		Name:    name,
		Comment: args.Comment,
		Context: Contexts[b[0]],
		Data:    b,
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	enc.Encode(output)
}
