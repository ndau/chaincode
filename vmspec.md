# The Virtual Machine

The virtual machine is a simple stack machine that is expected to evaluate an expression. It lacks any constructs that would enable looping (although it does have some aggregation operations).

Because we expect that vm scripts will end up as data on the blockchain, a compact representation of the VM is important. Consequently, the VM compiles to a bytecode representation.

Memory size is not as important, so the VM manipulates a typed stack, where the stack values can be any of a number of different datatypes.

When invoked, the VM's stack may already contain specific values in well-known locations. These values are part of the contract for each invocation (known as a context).

<!-- The function's return value is always the top item on the stack at exit (if multiple values should be returned, they must be returned in a list). The type and range of the return value is well-defined in advance; violation of these specifications is an error. Since errors have no way of being expressed further other than logging them, the function contract also defines the semantic interpretation of errors (for example, in some contexts, an error can be interpreted as a zero result). -->

## Data types

### Numbers
Numbers are 64-bit signed values; always integers. Overflow of a 64-bit value is an error.

### Timestamps
Timestamps are a 64-bit unsigned number of microseconds since the epoch. They can be added and subtracted but not multiplied or divided.

### Bytes
Bytes is the term for an arbitrary array of bytes. They can be compared only with other Bytes, and mathematical operations have no meaning.

### Lists
Lists are opaque references to a linear array of any of the datatypes. Duplicating the reference copies the list. Slices create new lists.

### Structs
Structs are opaque references to an object with fields in a specific order. Fields are identified by a field index which is always the following byte to the opcode (this limits the number of fields in a struct to 256). Fields also conform to the native data types.

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

Programs must conform to the following:
* Be no longer than MAX_LENGTH in bytes
* Every opcode must be defined within a function; functions must have matching `def` and `enddef` opcodes.
* Every conditional block must have a corresponding `endif` and at most one `else` opcode
* Nested conditionals must terminate before the termination of the parent clause
* The context identifies the expected type and order of the values on the stack, as well as intended return type. If anything fails to match with the runtime environment, the script fails.

Programs that fail to pass these validity checks will not be run and will be treated as if they errored on execution.

