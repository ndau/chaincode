# Chaincode

Chaincode is designed to present a decidedly NOT Turing-complete virtual machine for use in expressing formulas in a way that looks like data rather than code.

We have several places in the ndau system where it would be useful to be able to express a formula. We don't need or want the expressive power of a complete programming language. We want to be able to create small bits of code that are easily testable and not easily exploitable.

It is useful in several situations:

* Calculating the payouts to co-stakers for node operations
* Expressing the mechanism by which a node will decide how/when to generate EAI transactions
* Defining the formula for node quality that is evaluated to determine node ranking
* Expressing the mechanism by which market price is determined by combining price reports

This repository contains several pieces:

* A [spec for the virtual machine](vmspec.md).
* A library implementing the virtual machine itself.
* An assembler that can take text files corresponding to opcode-level instructions and create a properly formatted set of bytecodes.
* A compiler that describes a higher-level language for writing scripts.

