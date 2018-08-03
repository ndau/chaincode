
# Opcodes for Chaincode

This file is generated automatically; DO NOT EDIT.

## Implemented and Enabled Opcodes

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
0x00|Nop|No-op - has no effect.||
0x01|Drop|Discards the value on top of the stack.|A B|A
0x02|Drop2|Discards the top two values.|A B C|A
0x05|Dup|Duplicates the top of the stack.|A B|A B B
0x06|Dup2|Duplicates the top two items.|A B C|A B C B C
0x09|Swap|Exchanges the top two items on the stack.|A B C|A C B
0x0c|Over|Duplicates the second item on the stack to the top of the stack.|A B|A B A
0x0d|Pick|The item back in the stack by the specified offset is copied to the top.|A B C D|A B C D B
0x0e|Roll|The item back in the stack by the specified offset is moved to the top.|A B C D|A C D B
0x0f|Tuck|The top of the stack is dropped N entries back into the stack after removing it from the top.|A B C D|A D B C
0x10|Ret|Terminates the function; the values on the stack (if any) are the return values.||
0x11|Fail|Terminates the function and indicates an error.||
0x20|Zero (False)|Pushes 0 onto the stack.||0
0x21|Push1|Evaluates the next n bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x22|Push2|Evaluates the next 2 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x23|Push3|Evaluates the next 3 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x24|Push4|Evaluates the next 4 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x25|Push5|Evaluates the next 5 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x26|Push6|Evaluates the next 6 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x27|Push7|Evaluates the next 7 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x28|Push8|Evaluates the next 8 bytes as a signed little-endian numeric value and pushes it onto the stack.||A
0x29|PushB|Pushes the specified number of following bytes onto the stack as a Bytes object.||"ABC"
0x2a|One (True)|Pushes 1 onto the stack.||1
0x2b|Neg1|Pushes -1 onto the stack.||-1
0x2c|PushT|Concatenates the next 8 bytes and pushes them onto the stack as a timestamp.||timestamp A
0x2d|Now|Pushes the current timestamp onto the stack.||(current time as timestamp)
0x2e|PushA|Evaluates a to make sure it is formatted as a valid ndau-style address; if so, pushes it onto the stack as a Bytes object. If not, error.||nda234...4b3
0x2f|Rand|Pushes a 64-bit random number onto the stack. Note that 'random' may have special meaning depending on context; in particular, repeated uses of this opcode may (and most likely will) return the same value within a given context.||
0x30|PushL|Pushes an empty list onto the stack.||[]
0x40|Add|Adds the top two numeric values on the stack and puts their sum on top of the stack. attempting to add non-numeric values is an error.|A B|A+B
0x41|Sub|Subtracts the top numeric value on the stack from the second and puts the difference on top of the stack. attempting to subtract non-numeric values is an error.|A B|A-B
0x42|Mul|Multiplies the top two numeric values on the stack and puts their product on top of the stack. attempting to multiply non-numeric values is an error.|A B|A*B
0x43|Div|Divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero.|A B|int(A/B)
0x44|Mod|Divides the second numeric value on the stack by the top and puts the integer remainder on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero.|A B|A % B
0x45|DivMod|Divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack and the remainder in the second item on the stack. attempting to divide non-numeric values is an error, as is dividing by zero.|A B|A%B int(A/B)
0x46|MulDiv|Multiplies the third numeric item on the stack by the fraction created by dividing the second numeric item by the top; guaranteed not to overflow as long as the fraction is less than 1. An overflow is an error.|A B C|int(A*(B/C))
0x48|Not|If the top of the stack is 0, it is replaced by 1 -- otherwise, it is replaced by 0.|5 6 7|5 6 0
0x49|Neg|The sign of the number on top of the stack is negated.|A|-A
0x4a|Inc|Adds 1 to the number on top of the stack.|A|A+1
0x4b|Dec|Subtracts 1 from the number on top of the stack.|A|A-1
0x4d|Eq|Compares (and discards) the two top stack elements. If they are equal in both type and value, leaves TRUE (1) on top of the stack, otherwise leaves FALSE (0) on top of the stack.|A B|FALSE
0x4e|Gt|Compares (and discards) the two top stack elements. If the types are different, fails execution. If the types are the same, compares the values, and leaves TRUE when: a) the top Number or Timestamp or ID is numerically greater than the second, b) the top list is longer than the second list, c) the top struct compares greater than the second by iterating in field order and using rules a) and b). .|A B|TRUE
0x4f|Lt|Like gt, using less than instead of greater than.|A B|FALSE
0x50|Index|Selects a zero-indexed element (the index is the top of the stack) from a list reference which is the second item on the stack (both are discarded) and leaves it on top of the stack. Error if index is out of bounds or a list is not on top of the stack.|[X Y Z] 2|Z
0x51|Len|Returns the length of a list.|[X Y Z]|3
0x52|Append|Creates a new list, appending the new value to it.|[X Y] Z|[X Y Z]
0x53|Extend|Generates a new list by concatenating two other lists.|[X Y] [Z]|[X Y Z]
0x54|Slice|Expects a list and two indices on top of the stack. Creates a new list containing the designated subset of the elements in the original slice.|[X Y Z] 1 3|[Y Z]
0x60|Field|Retrieves a field at index f from a struct; if the index is out of bounds, fails.|X|X.f
0x70|FieldL|Makes a new list by retrieving a given field from all of the structs in a list.|[X Y Z]|[X.f Y.f Z.f]
0x80|Def|Defines function block n, where n is a number larger than any previously defined function in this script. Function 0 is called by the system. Every function must be terminated by end, and function definitions may not be nested.||
0x81|Call|Calls the function block, provided that its ID is greater than the index of the function block currently executing (recursion is not permitted). The function runs with a new stack which is initialized with the top n values of the current stack (which are copied, NOT popped). Upon return, the top value on the function's stack is pushed onto the caller's stack.||
0x82|Deco|Decorates a list of structs (on top of the stack, which it pops) by applying the function block to each member of the struct, copying n stack entries to the function block's stack, then copying the struct itself; on return, the top value of the function block stack is appended to the list entry. The resulting new list is pushed onto the stack.||
0x88|EndDef|Ends a function definition; always required.||
0x89|IfZ|If the top stack item is zero, executes subsequent code. The top stack item is discarded.||
0x8a|IfNZ|If the top stack item is nonzero, executes subsequent code. The top stack item is discarded.||
0x8e|Else|If the code immediately following an if was not executed, this code (up to end) will be; otherwise it will be skipped.||
0x8f|EndIf|Terminates a conditional block; if this opcode is missing for any block, the program is invalid.||
0x90|Sum|Given a list of numbers, sums all the values in the list.|[2 12 4]|18
0x91|Avg|Given a list of numbers, averages all the values in the list.|[2 12 4]|6
0x92|Max|Given a list of numbers, finds the maximum value.|[2 12 4]|12
0x93|Min|Given a list of numbers, finds the minimum value.|[2 12 4]|2
0x94|Choice|Selects an item at random from a list and leaves it on the stack as a replacement for the list.|[X Y Z]|
0x95|WChoice|Selects an item from a list of structs weighted by the given field index.|[X Y Z] f|
0x96|Sort|Sorts a list of structs by a given field.|[X Y Z] f|The list sorted by field f
0x97|Lookup|Selects an item from a list of structs by applying the function block to each item in order, copying n stack entries to the function block's stack, then copying the struct itself; returns the index of the first item in the list where the result is a nonzero number; throws an error if no item returns a nonzero number.|[X Y Z]|i
# Disabled Opcodes

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
||There are no disabled opcodes at the moment.||
