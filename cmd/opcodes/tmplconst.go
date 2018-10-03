package main

// we expect this to be invoked on OpcodeData
const tmplConstDef = `
// Code generated automatically by "make generate"; DO NOT EDIT.

package main

// Predefined constants available to chasm programs.

func predefinedConstants() map[string]string {
	k := map[string]string{
		"EVENT_DEFAULT":                "0",
		"EVENT_TRANSFER":               "1",
		"EVENT_CHANGETRANSFERKEY":      "2",
		"EVENT_RELEASEFROMENDOWMENT":   "3",
		"EVENT_CHANGESETTLEMENTPERIOD": "4",
		"EVENT_DELEGATE":               "5",
		"EVENT_COMPUTEEAI":             "6",
		"EVENT_LOCK":                   "7",
		"EVENT_NOTIFY":                 "8",
		"EVENT_SETREWARDSTARGET":       "9",
		"EVENT_CLAIMACCOUNT":           "10",
		"EVENT_TRANSFERANDLOCK":        "11",
		"EVENT_GTVALIDATORCHANGE":      "0XFF",
{{range . -}}
		"{{.Name}}": "{{printf "%d" .Value}}",
{{end}}
	}
	return k
}
`
