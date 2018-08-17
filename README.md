# Chaincode

Chaincode is designed to present a decidedly NOT Turing-complete virtual machine for use in expressing formulas in a way that looks like data rather than code.

We have several places in the ndau system where it would be useful to be able to express a formula. We don't need or want the expressive power of a complete programming language. We want to be able to create small bits of code that are easily testable and not easily exploitable.

It is useful in several situations:

* Calculating the payouts to co-stakers for node operations
* Expressing the mechanism by which a node will decide how/when to generate EAI transactions
* Defining the formula for node quality that is evaluated to determine node ranking
* Expressing the mechanism by which market price is determined by combining price reports

## Contents
This repository contains several pieces:

* A [spec for the virtual machine](vmspec.md).
* A library (pkg/vm) implementing the virtual machine itself.
* An assembler (chasm) that can take text files corresponding to opcode-level instructions and create a properly formatted set of bytecodes.
* A code generator (cmd/opcodes) that contains the definitions of all the opcodes (see cmd/chasm/opcodedata.go) and uses that to generate:
    * Documentation
    * Keywords for chasm
    * Keywords for the mini-assembler in the VM
    * Table of enabled opcodes (allowing us to disable opcodes on the fly if necessary)
    * Table of bytecounts for multibyte instructions

## Building the code generator
To build the code generator and get all the files up-to-date:

```sh
cd cmd/opcodes
go get
cd ../..
make generate
```

## Testing the VM library

You should run the code generator first to me sure everything is up to date. There is no explicit dependency for make fuzz because we want to be able to run it over and over.

```sh
make fuzz
```

It first runs the normal VM tests with the -short flag (which excludes the fuzz tests), then runs a set of fuzz tests that are fairly comprehensive but designed to be CI-friendly (they complete within 30 seconds).

### Other options

If you run tests with `go test` without the -short flag, you'll also get a short run of the fuzz tests that are designed to complete in less than 10 seconds.

If you want to really exercise things, you can run `make fuzzmillion`, which runs each of the fuzz tests one million times.

## Building chasm

You need to have [pigeon](https://github.com/mna/pigeon) installed in order to build the parser.

```sh
go get -u github.com/mna/pigeon
```

To build and test:
```sh
cd cmd/chasm
make test
```

To use the assembler to compile examples:
```sh
make examples
```

Running the assembler with `chasm -h` will print a help message.

