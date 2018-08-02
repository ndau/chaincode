package main

// we expect this to be invoked on OpcodeData
const tmplOpcodesMiniAsm = `
package vm

// these are the opcodes supported by mini-asm
var opcodeMap = map[string]Opcode{
{{range .MiniAsm -}}
	"{{tolower .Name}}": Op{{.Name}},
{{end}}
}
`
