package vm

// ----- ---- --- -- -
// Copyright 2019, 2020 The Axiom Foundation. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----


import (
	"encoding/json"
	"io"
)

// ChasmBinary defines the "binary" (assembled) format of a vm
type ChasmBinary struct {
	Name    string   `json:"name"`
	Comment string   `json:"comment"`
	Data    []Opcode `json:"data"`
}

// Serialize takes a stream of bytes and sends it to
// a Writer in ChasmBinary format
func Serialize(name string, comment string, b []byte, w io.Writer) error {
	opcodes := make([]Opcode, len(b))
	for i := range b {
		opcodes[i] = Opcode(b[i])
	}
	output := ChasmBinary{
		Name:    name,
		Comment: comment,
		Data:    opcodes,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

// Deserialize takes a reader and extracts a ChasmBinary from it
func Deserialize(r io.Reader) (ChasmBinary, error) {
	var input ChasmBinary
	err := json.NewDecoder(r).Decode(&input)
	return input, err
}
