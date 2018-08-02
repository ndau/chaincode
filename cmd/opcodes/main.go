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

func doOpcodeDoc(tname string, w io.Writer) error {
	var tmpl = template.Must(template.New(tname).Parse(tmplOpcodeDoc))
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
	funcMap := template.FuncMap{
		"tolower": strings.ToLower,
	}

	var tmpl = template.Must(template.New(tname).Funcs(funcMap).Parse(ts))

	return tmpl.Execute(w, opcodeData)
}

func main() {
	var args struct {
		Opcodes string `arg:"-o" help:"opcodes doc file -- ./opcodes.md"`
		Defs    string `arg:"-d" help:"opcode definition file -- ./pkg/vm/opcodes.go"`
		MiniAsm string `arg:"-m" help:"mini-assembler opcodes -- ./pkg/vm/miniasmOpcodes.go"`
		Pigeon  string `arg:"-p" help:"pigeon grammar for opcodes -- ./cmd/chasm/chasm.peggo (modifies this file)"`
	}
	arg.MustParse(&args)

	var err error

	if args.Opcodes != "" {
		f := os.Stdout
		if args.Opcodes != "-" {
			f, err = os.Create(args.Opcodes)
			defer f.Close()
			if err != nil {
				panic(err)
			}
		}
		doOpcodeDoc(args.Opcodes, f)
	}

	if args.Defs != "" {
		f := os.Stdout
		ondisk := false
		if args.Defs != "-" {
			ondisk = true
			f, err = os.Create(args.Defs)
			if err != nil {
				panic(err)
			}
		}

		doOpcodesGo(args.Defs, tmplOpcodesDef, f)
		if ondisk {
			f.Close()
			gofmtFile(args.Defs)
		}
	}

	if args.MiniAsm != "" {
		f := os.Stdout
		ondisk := false
		if args.MiniAsm != "-" {
			ondisk = true
			f, err = os.Create(args.MiniAsm)
			if err != nil {
				panic(err)
			}
		}

		doOpcodesGo(args.MiniAsm, tmplOpcodesMiniAsm, f)
		if ondisk {
			f.Close()
			gofmtFile(args.MiniAsm)
		}
	}
}
