package main

// we expect this to be invoked on OpcodeData
const tmplOpcodesMiniAsm = `
// This file is generated automatically; DO NOT EDIT.

package vm

// these are the opcodes supported by mini-asm
var opcodeMap = map[string]Opcode{
{{range .EnabledWithSynonyms -}}
	"{{tolower .Name}}": Op{{.Name}},
{{end}}
}
`
