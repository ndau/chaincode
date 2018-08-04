package main

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	arg "github.com/alexflint/go-arg"
)

// This command generates a number of files in the chaincode project from the opcode data.
// By default, it needs no parameters -- it just generates all of the files at once
// from the information in the opcodedata.go file.

var funcMap = template.FuncMap{
	"tolower": strings.ToLower,
	"getparm": getParm,
	"nbytes":  nbytes,
}

func doOpcodeDoc(tname string, ts string, w io.Writer) error {
	var tmpl = template.Must(template.New(tname).Funcs(funcMap).Parse(ts))
	err := tmpl.Execute(w, opcodeData)
	if err != nil {
		return err
	}
	return nil
}

func gofmtFile(name string) error {
	cmd := exec.Command("gofmt", "-w", name)
	return cmd.Run()
}

func doOpcodesGo(tname string, ts string, w io.Writer) error {
	var tmpl = template.Must(template.New(tname).Funcs(funcMap).Parse(ts))

	return tmpl.Execute(w, opcodeData)
}

func generateGoFile(name, tmpl string) {
	f := os.Stdout
	var err error
	ondisk := false
	if name != "-" {
		ondisk = true
		f, err = os.Create(name)
		if err != nil {
			panic(err)
		}
	}

	err = doOpcodesGo(name, tmpl, f)
	if err != nil {
		panic(err)
	}
	if ondisk {
		f.Close()
		gofmtFile(name)
	}
}

func main() {
	var args struct {
		Opcodes string `arg:"-o" help:"opcodes doc file -- ./opcodes.md"`
		Defs    string `arg:"-d" help:"opcode definition file -- ./pkg/vm/opcodes.go"`
		MiniAsm string `arg:"-m" help:"mini-assembler opcodes -- ./pkg/vm/miniasmOpcodes.go"`
		Extra   string `arg:"-e" help:"extrabytes helper for opcodes -- ./pkg/vm/extrabytes.go"`
		Pigeon  string `arg:"-p" help:"pigeon grammar for opcodes -- ./cmd/chasm/chasm.peggo (modifies this file)"`
	}
	arg.MustParse(&args)

	var err error

	if args.Defs != "" {
		generateGoFile(args.Defs, tmplOpcodesDef)
	}

	if args.MiniAsm != "" {
		generateGoFile(args.MiniAsm, tmplOpcodesMiniAsm)
	}

	if args.Extra != "" {
		generateGoFile(args.Extra, tmplOpcodesExtra)
	}

	if args.Opcodes != "" {
		f := os.Stdout
		if args.Opcodes != "-" {
			f, err = os.Create(args.Opcodes)
			defer f.Close()
			if err != nil {
				panic(err)
			}
		}
		err = doOpcodeDoc(args.Opcodes, tmplOpcodeDoc, f)
		if err != nil {
			panic(err)
		}
	}

	if args.Pigeon != "" {
		var w io.WriteCloser = os.Stdout
		if args.Pigeon != "-" {
			w, err = NewInjectionWriter(args.Pigeon, "// VVVVV---GENERATED", "// ^^^^^---GENERATED")
			if err != nil {
				panic(err)
			}
			defer w.Close()
		}
		err = doOpcodeDoc(args.Pigeon, tmplOpcodesPigeon, w)
		if err != nil {
			panic(err)
		}
	}

}
