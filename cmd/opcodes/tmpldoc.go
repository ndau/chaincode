package main

// we expect this to be invoked on OpcodeData
const tmplOpcodeDoc = `
# Opcodes for Chaincode

This file is generated automatically; DO NOT EDIT.

## Implemented and Enabled Opcodes

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
{{range .Enabled -}}
    {{- printf "0x%02x" .Value}}|
    {{- .Name}}{{if .Synonym}} ({{.Synonym}}){{end}}|
    {{- .Summary}}|
    {{- .Example.Pre}}|
    {{- .Example.Inst}}|
    {{- .Example.Post}}
{{end -}}

# Disabled Opcodes

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
{{range .Disabled -}}
    {{- printf "0x%02x" .Value}}|
    {{- .Name}}|
    {{- .Summary}}|
    {{- .Example.Pre}}|
    {{- .Example.Inst}}|
    {{- .Example.Post}}
{{else -}}
||There are no disabled opcodes at the moment.||
{{end -}}

`
