# The Virtual Machine

The virtual machine is a simple stack machine that is expected to evaluate an expression. It lacks any constructs that would enable looping (although it does have some aggregation operations).

It is invoked with a stack that contains specific values in well-known locations. These values are part of the contract for each invocation.

The function's return value is always the top item on the stack at exit (if multiple values should be returned, they must be returned in a list). The type and range of the return value is well-defined in advance; violation of these specifications is an error. Since errors have no way of being expressed further other than logging them, the function contract also defines the semantic interpretation of errors (for example, in some contexts, an error can be interpreted as a zero result).

## Data types

### Numbers
Numbers are 64-bit values; always integers. Overflow of a 64-bit value is an error.

### Timestamps
Timestamps are a 64-bit number of microseconds since the epoch. They can be added and subtracted but not multiplied or divided.

### Structs
Structs are opaque references to an object with fields in a specific order. Fields are identified by a field index which is the following byte. Fields also conform to the native data types.

Structs cannot be created or modified by functions.

### Lists
Lists are opaque references to a linear array of any of the datatypes. All items in a list must be the same type. Duplicating the reference does not copy the list. Slices create new lists.

## Script Structure
Every script has a one-byte preamble defining its context. The context implies:

* The type of its return value
* The quantity and types of values it expects on the stack at entry

Contexts will be published; a change in any of these parameters will require creating a new context.

The remainder of the bytes in a script are the opcodes.

### Opcodes
(see opcodes file)

## Validity

Before being executed, programs are examined.

Programs must conform to the following:
* Be no longer than MAX_LENGTH in bytes
* Every block must have a corresponding end and at most one else opcode
* Nested ifs must terminate before the termination of the parent clause
* The context identifies the expected type and order of the values on the stack, as well as intended return type. If anything fails to match with the runtime environment, the script fails.

Programs that fail to pass these validity checks will not be run and will be treated as if they errored on execution.

