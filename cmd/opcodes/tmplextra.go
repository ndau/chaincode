package main

// we expect this to be invoked on OpcodeData
const tmplOpcodesExtra = `
// Code generated automatically by "make generate" -- DO NOT EDIT

package vm

// extraBytes returns the number of extra bytes associated with a given opcode
func extraBytes(code Chaincode, offset int) int {
	numExtra := 0
	op := code[offset]
	switch op {
{{- range .Enabled -}}{{if not (eq (len .Parms) 0)}}
	case Op{{.Name}}:
		numExtra = {{nbytes .}}
{{- end}}{{end}}
	}
	return numExtra
}
`
