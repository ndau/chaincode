package main

import "fmt"

// Value: 0,
// Name: "nop",
// Summary: "no-op - has no effect",
// Doc: "",
// Example: example{
//     Pre: "A B",
//     Inst: "nop 3",
//     Post: "A B",
// },
// Parms: []parm{indexParm{}, shiftParm{}},
// ErrorNotes: "",

type example struct {
	Pre  string
	Inst string
	Post string
}

type parm interface {
	Nbytes() int
	String() string
}

type stackOffsetParm struct{}

func (p stackOffsetParm) Nbytes() int {
	return 1
}

func (p stackOffsetParm) String() string {
	return "n"
}

type countParm struct{}

func (p countParm) Nbytes() int {
	return 1
}

func (p countParm) String() string {
	return "n"
}

type functionIDParm struct{}

func (p functionIDParm) Nbytes() int {
	return 1
}

func (p functionIDParm) String() string {
	return "id"
}

type indexParm struct {
	Name string
}

func (p indexParm) Nbytes() int {
	return 1
}

func (p indexParm) String() string {
	return p.Name
}

type dataParm struct {
	N int
}

func (p dataParm) Nbytes() int {
	return p.N
}

func (p dataParm) String() string {
	return fmt.Sprintf("(%d data bytes)", p.N)
}

type opcodeInfo struct {
	Value   byte
	Name    string
	Synonym string
	Summary string
	Doc     string
	Errors  string
	Example example
	Parms   []parm
	Enabled bool
}

type opcodeInfos []opcodeInfo

// selectEnabled returns the subset of the opcodeInfo that matches
// the state of the enabled flag. If the withSynonym flag is specified,
// it also generates records for any synonyms
func (o opcodeInfos) selectEnabled(enabled bool, withSynonym bool) opcodeInfos {
	o2 := make(opcodeInfos, 0)
	for i := range o {
		if o[i].Enabled == enabled {
			o2 = append(o2, o[i])
			if withSynonym && o[i].Synonym != "" {
				syn := o[i]
				syn.Name = syn.Synonym
				o2 = append(o2, syn)
			}
		}
	}
	return o2
}

func (o opcodeInfos) Enabled() opcodeInfos {
	return o.selectEnabled(true, false)
}

func (o opcodeInfos) EnabledWithSynonyms() opcodeInfos {
	return o.selectEnabled(true, true)
}

func (o opcodeInfos) Disabled() opcodeInfos {
	return o.selectEnabled(false, false)
}
