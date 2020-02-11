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
	"encoding"
	"encoding/base64"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

// Chaincode is the type for the VM bytecode program
type Chaincode []Opcode

// ConvertToOpcodes accepts an array of bytes and returns a Chaincode (array of opcodes)
func ConvertToOpcodes(b []byte) Chaincode {
	ops := make([]Opcode, len(b))
	for i := range b {
		ops[i] = Opcode(b[i])
	}
	return Chaincode(ops)
}

// ToChaincode converts a byte slice into Chaincode
func ToChaincode(b []byte) Chaincode {
	return ConvertToOpcodes(b)
}

// Bytes converts chaincode into a byte slice
func (c Chaincode) Bytes() []byte {
	bytes := make([]byte, len(c))
	for i := range c {
		bytes[i] = byte(c[i])
	}
	return bytes
}

// IsValid tests if an array of opcodes is a potentially valid
// Chaincode program.
func (c Chaincode) IsValid() error {
	if c == nil {
		return ValidationError{"missing code"}
	}
	if len(c) > maxTotalLength {
		return ValidationError{"code and data combined are too long"}
	}
	// make sure the executable part of the code is valid
	_, _, err := validateStructure(c)
	if err != nil {
		return err
	}

	if c.CodeSize() > maxCodeLength {
		return ValidationError{"code is too long"}
	}

	// now generate a bitset of used opcodes from the instructions
	usedOpcodes := getUsedOpcodes(generateInstructions(c))
	// if it's not a proper subset of the enabled opcodes, don't let it run
	if !usedOpcodes.IsSubsetOf(EnabledOpcodes) {
		return ValidationError{"code contains illegal opcodes"}
	}
	return nil
}

var _ encoding.TextMarshaler = (*Chaincode)(nil)
var _ encoding.TextUnmarshaler = (*Chaincode)(nil)

// MarshalText implements encoding.TextMarshaler
func (c Chaincode) MarshalText() (text []byte, err error) {
	base64.StdEncoding.Encode(text, c.Bytes())
	return
}

// UnmarshalText implements encoding.TextUnmarshaler
func (c *Chaincode) UnmarshalText(text []byte) error {
	bytes := make([]byte, 0, base64.StdEncoding.DecodedLen(len(text)))
	_, err := base64.StdEncoding.Decode(bytes, text)
	if err != nil {
		return err
	}
	*c = ToChaincode(bytes)
	return nil
}

var _ msgp.Marshaler = (*Chaincode)(nil)
var _ msgp.Unmarshaler = (*Chaincode)(nil)
