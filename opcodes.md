# Implemented Opcodes for Chaincode

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
0x00|nop|no-op - has no effect||
0x01|drop|discards the value on top of the stack|A B|A
0x02|drop2|discards the top two values|A B C|A
0x05|dup|duplicates the top of the stack|A B|A B B
0x06|dup2|duplicates the top two items|A B C|A B C B C
0x09|swap|exchanges the top two items on the stack|A B C|A C B
0x0C|over|duplicates the second item on the stack to the top of the stack|A B|A B A
0x0D|pick n|the item N back in the stack is copied to the top; pick 0 is dup; pick 1 is over|A B C D|A B C D B (if n is 2)
0x0E|roll n|the item N back in the stack is moved to the top; roll 0 is nop, roll 1 is swap|A B C D|A C D B (if n is 2)
0x0F|tuck n|the top of the stack is inserted N back, removing it from the top. tuck 0 is nop, tuck 1 is swap|A B C D|A D B C (if n is 2)
0x10|ret|terminates the function; the values on the stack (if any) are the return values.||
0x11|fail|terminates the function and indicates an error||
0x20|zero, false|Pushes 0 onto the stack||0
0x21-0x28|pushN (where N is 1-8)|evaluates the next n bytes as a signed little-endian numeric value and pushes it onto the stack||A
0x29|pushb n|Pushes n following bytes onto the stack as a Bytes object||
0x2A|one, true|Pushes 1 onto the stack||1
0x2B|neg1|Pushes -1 onto the stack||-1
0x2C|pusht|concatenates the next 8 bytes and pushes them onto the stack as a timestamp||timestamp A
0x2D|now|Pushes the current timestamp onto the stack. Note that "current" may have special meaning depending on the context; in particular, repeated uses of this opcode may (and most likely will) return the same value within a given context.||
0x2E|pushaddr a|Evaluates a to make sure it is a valid ndau-style address; if so, pushes it onto the stack as a Bytes object. If not, error.||
0x2F|rand|Pushes a 64-bit random number onto the stack. Note that "random" may have special meaning depending on context; in particular, repeated uses of this opcode may (and most likely will) return the same value within a given context.||
0x30|pushl|Pushes an empty list onto the stack||
0x40|add|adds the top two numeric values on the stack and puts their sum on top of the stack. attempting to add non-numeric values is an error|A B|A+B
0x41|sub|subtracts the top numeric value on the stack from the second and puts the difference on top of the stack. attempting to subtract non-numeric values is an error|A B|A-B
0x42|mul|multiplies the top two numeric values on the stack and puts their product on top of the stack. attempting to multiply non-numeric values is an error|A B|A*B
0x43|div|divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero|A B|int(A/B)
0x44|mod|divides the second numeric value on the stack by the top and puts the integer remainder on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero|A B|A % B
0x45|divmod|divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack and the remainder in the second item on the stack. attempting to divide non-numeric values is an error, as is dividing by zero|A B|A%B int(A/B)
0x46|muldiv|multiplies the third numeric item on the stack by the fraction created by dividing the second numeric item by the top; guaranteed not to overflow as long as the fraction is less than 1. An overflow is an error.|A B C|int(A*(B/C))
0x48|not|if the top of the stack is 0, it is replaced by 1 -- otherwise, it is replaced by 0||
0x49|neg|the sign of the number on top of the stack is negated|A|-A
0x4A|inc|adds 1 to the number on top of the stack|A|A+1
0x4B|dec|subtracts 1 from the number on top of the stack|A|A-1
0x4D|eq|Compares (and discards) the two top stack elements. If they are equal in both type and value, leaves TRUE (1) on top of the stack, otherwise leaves FALSE (0) on top of the stack.|A B|FALSE
0x4E|gt|Compares (and discards) the two top stack elements. If the types are different, fails execution. If the types are the same, compares the values, and leaves TRUE when: a) the top Number or Timestamp or ID is numerically greater than the second, b) the top list is longer than the second list, c) the top struct compares greater than the second by iterating in field order and using rules a) and b). |A B|TRUE
0x4F|lt|like gt, using less than instead of greater than|A B|FALSE
0x50|index|selects a zero-indexed element (the index is the top of the stack) from a list reference which is the second item on the stack (both are discarded) and leaves it on top of the stack. Error if index is out of bounds or a list is not on top of the stack.|[X Y Z] 2|Z
0x51|len|Returns the length of a list|[X Y Z]|3
0x52|append|creates a new list, appending the new value to it|[X Y] Z|[X Y Z]
0x53|extend|generates a new list by concatenating two other lists|[X Y] [Z]|[X Y Z]
0x54|slice|Expects a list and two indices on top of the stack. Creates a new list containing the designated subset of the elements in the original slice.|[X Y Z] 1 3|[Y Z]
0x60|field f|retrieves a field at index f from a struct; if the index is out of bounds, fails|X|X.f
0x70|fieldl f|makes a new list by retrieving a given field from all of the structs in a list|[X Y Z]|[X.f Y.f Z.f]
0x80|def n|defines function block n, where n is a number larger than any previously defined function in this script. Function 0 is called by the system. Every function must be terminated by end, and function definitions may not be nested.||
0x81|call n m|calls function block n, provided that n is greater than the index of the function block currently executing (recursion is not permitted). The function runs with a new stack which is initialized with the top m values of the current stack (which are copied, NOT popped). Upon return, the top value on the function's stack is pushed onto the caller's stack. Functions must return a single Value (it may be a List or Struct).||
0x82|deco n m|Decorates a list of structs (on top of the stack, which it pops) by applying function block n to each member of the struct, copying m stack entries to the function block's stack, then copying the struct itself; on return, the top value of the function block stack is appended to the list entry. The resulting new list is pushed onto the stack.||
0x88|enddef|Ends a function definition; always required||
0x89|ifz|if the top stack item is zero, executes subsequent code. The top stack item is discarded||
0x8A|ifnz|if the top stack item is nonzero, executes subsequent code. The top stack item is discarded||
0x8E|else|if the code immediately following an if was not executed, this code (up to end) will be; otherwise it will be skipped||
0x8F|endif|terminates a conditional block; if this opcode is missing for any block, the program is invalid.||
0x90|sum|Given a list of numbers, sums all the values in the list.|[2 12 4]|18
0x91|avg|Given a list of numbers, averages all the values in the list.|[2 12 4]|6
0x92|max|Given a list of numbers, finds the maximum value|[2 12 4]|12
0x93|min|Given a list of numbers, finds the minimum value|[2 12 4]|2
0x94|choice|selects an item at random from a list and leaves it on the stack as a replacement for the list|[X Y Z]|
0x95|wchoice f|selects an item from a list of structs weighted by the given field index|[X Y Z] f|
0x96|sort f|sorts a list of structs by a given field|[X Y Z] f|The list sorted by field f
0x97|lookup n m|selects an item from a list of structs by applying function block n to each item in order, copying m stack entries to the function block's stack, then copying the struct itself; returns the index of the first item in the list where the result is a nonzero number; throws an error if no item returns a nonzero number|[X Y Z]|i

# Opcodes under consideration

Value|Opcode|Meaning|Stack before|Stack after
----|----|----|----|----
0x61|set f|sets a field at index f in a struct; if the index is out of bounds, fails (is it a good idea to have set or should system structs be readonly?)|X v|X.f <- v
0xxx|intersect|Given two lists on top of the stack, pops both and returns a new list containing only the items in both lists that compare as equal. The ordering of the result list is not guaranteed.|[A B C D] [A C E]|[A C]
0xxx|verify|Top of stack: signature (as Bytes); #2 on stack: message hash (as Bytes); #3: public key (as Bytes); Verifies that the signature is a valid signature of the message hash with the given public key. (This is done by the tx prevalidation)||
0xxx|hash|Calculates a hash of the given block of bytes (why would we need this?)||
0xxx|addr|Generates an address, given a key (as bytes) (why would we need this?)|Key|Addr
0xxx|filter n m|calls n on every element of list and returns a list containing only the elements for which n was true (is it actually useful?)|

bitwise: we might want to consider supporting bitwise operations so people can do flag testing
    AND
    OR
    XOR
    NOT

list: some more list operations may be interesting, but could also be dangerous
    ALL
    ANY
    COUNT
    (reduce?)

