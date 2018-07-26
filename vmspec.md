# Chaincode

We have several places in the ndau system where it would be useful to be able to express a formula. Chaincode provides a virtual machine and an assembler (called chasm) for that VM.

It is almost certainly unwise to provide a Turing-complete programming language, but a mechanism for expressing a formula that does not permit iteration is useful in several situations:

* Expressing rules for the validity of transactions on accounts (M of N, 1 + M of N, spending limits, etc)
* Calculating the payouts to co-stakers for node operations
* Expressing the mechanism by which a node will decide how/when to generate EAI transactions
* Defining the formula for node quality that is evaluated to determine node ranking
* Expressing the mechanism by which market price is determined by combining price reports

# The Virtual Machine

The virtual machine is a stack machine that is expected to evaluate an expression. It lacks any constructs that would enable looping (although it does have some aggregation operations).

Because we expect that vm scripts will end up as data on the blockchain, a compact representation of the VM is important. Consequently, the VM compiles to a bytecode representation.

Memory size is not as important, so the VM manipulates a typed stack, where the stack values can be any of a number of different datatypes.

The VM operates through manipulation of a stack. Unlike hardware-based VMs, the values on the stack are typed and operations on the stack are defined with respect to specific types. (For example, to put a string on the stack, you push a string, not a sequence of embedded bytes.)

Every VM is executed in a specific context, and the context is defined when coding the VM. The context expresses the contract for what values will be on the stack when the VM is invoked, and what it is expected to return. All valid contexts are expected to be known in advance).

The function's "return value" is always the top item(s) on its stack at exit. The type and range of the return value is well-defined in advance; violation of these specifications is an error. Since errors have no way of being expressed further other than logging them (they run in the context of a node, not in a user environment), the function context also defines the semantic interpretation of errors (for example, in some contexts, an error can be interpreted as a zero result).


## Data types

### Numbers
Numbers are 64-bit signed values; always integers. Overflow of a 64-bit value is an error -- but the muldiv and divmod operations provide temporary 128-bit math features to avoid overflow.

### Timestamps
Timestamps are a 64-bit unsigned number of microseconds since the epoch. They can be added and subtracted but not multiplied or divided.

### Bytes
Bytes is the term for an arbitrary array of bytes. They can be compared only with other Bytes, and mathematical operations have no meaning.

### Lists
Lists are a linear array of any of the datatypes. Some operations on lists may require that all items in a list are of the same type. Lists contain copies of items, not references to them (duplicating a list and then modifying the contents of one of them will not modify the other one).

### Structs
Structs are opaque references to an object with fields in a specific order. Fields are identified by a field index which is always the following byte to the opcode (this limits the number of fields in a struct to 256). Fields also conform to the native data types.

Contexts define a set of constants that allow structs to be manipulated by name.

Structs cannot be created by the VM, only by its callers. However, they can be modified in a limited way by the `deco` opcode.

## Script Structure
Every script has a one-byte preamble defining its context. The context implies:

* The quantity and types of values it expects on the stack at entry
* The quantity and types of values it expects on the stack at exit
* How an error exit will be interpreted

Contexts are published in [contexts.md](contexts.md); a change in any of these parameters will require creating a new context.

The remainder of the bytes in a script are the opcodes.

All opcodes must be defined within functions defined using the `def` and `enddef` opcode pair. The parameter to `def` is the function index, and in a given script the functions must be defined in sequential numeric order, starting with 0. The zero function is considered "main" and execution of a script starts there. Functions are called using the `call` opcode, which has the restriction that it can only call functions whose number is strictly greater than the currently executing function (this constraint prevents recursion).

### Opcodes
(see [opcodes.md](opcodes.md))

## Validity

Before being executed, programs are examined.

Programs must conform to a variety of rules, including:
* Be no longer than MAX_LENGTH in bytes
* Every opcode must be defined within a function; functions must have matching `def` and `enddef` opcodes.
* Every conditional block must have a corresponding `endif` and at most one `else` opcode
* Nested conditionals must terminate before the termination of the parent clause
* The context identifies the expected type and order of the values on the stack, as well as intended return type. If anything fails to match with the runtime environment, the script fails.

Programs that fail to pass these validity checks will not be run and will be treated as if they errored on execution.

