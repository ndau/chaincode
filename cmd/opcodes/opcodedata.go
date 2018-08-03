package main

var opcodeData = opcodeInfos{
	opcodeInfo{
		Value:   0x00,
		Name:    "Nop",
		Summary: "No-op - has no effect.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "nop",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x01,
		Name:    "Drop",
		Summary: "Discards the value on top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "drop",
			Post: "A"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x02,
		Name:    "Drop2",
		Summary: "Discards the top two values.",
		Doc:     "",
		Example: example{
			Pre:  "A B C",
			Inst: "drop2",
			Post: "A"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x05,
		Name:    "Dup",
		Summary: "Duplicates the top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "dup",
			Post: "A B B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x06,
		Name:    "Dup2",
		Summary: "Duplicates the top two items.",
		Doc:     "",
		Example: example{
			Pre:  "A B C",
			Inst: "dup2",
			Post: "A B C B C"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x09,
		Name:    "Swap",
		Summary: "Exchanges the top two items on the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A B C",
			Inst: "swap",
			Post: "A C B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x0C,
		Name:    "Over",
		Summary: "Duplicates the second item on the stack to the top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "over",
			Post: "A B A"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x0D,
		Name:    "Pick",
		Summary: "The item back in the stack by the specified offset is copied to the top.",
		Doc:     "Pick 0 is the same as dup; pick 1 is over.",
		Example: example{
			Pre:  "A B C D",
			Inst: "pick 2",
			Post: "A B C D B"},
		Parms:   []parm{stackOffsetParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x0E,
		Name:    "Roll",
		Summary: "The item back in the stack by the specified offset is moved to the top.",
		Doc:     "Roll 0 is the same as nop, roll 1 is swap.",
		Example: example{
			Pre:  "A B C D",
			Inst: "roll 2",
			Post: "A C D B"},
		Parms:   []parm{stackOffsetParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x0F,
		Name:    "Tuck",
		Summary: "The top of the stack is dropped N entries back into the stack after removing it from the top.",
		Doc:     "Tuck 0 is the same as nop, tuck 1 is swap.",
		Example: example{
			Pre:  "A B C D",
			Inst: "tuck 2",
			Post: "A D B C"},
		Parms:   []parm{stackOffsetParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x10,
		Name:    "Ret",
		Summary: "Terminates the function; the values on the stack (if any) are the return values.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "ret",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x11,
		Name:    "Fail",
		Summary: "Terminates the function and indicates an error.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "fail",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x20,
		Name:    "Zero",
		Synonym: "False",
		Summary: "Pushes 0 onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "zero",
			Post: "0"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x21,
		Name:    "Push1",
		Summary: "Evaluates the next n bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push1",
			Post: "A"},
		Parms:   []parm{dataParm{1}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x22,
		Name:    "Push2",
		Summary: "Evaluates the next 2 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push2",
			Post: "A"},
		Parms:   []parm{dataParm{2}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x23,
		Name:    "Push3",
		Summary: "Evaluates the next 3 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push3",
			Post: "A"},
		Parms:   []parm{dataParm{3}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x24,
		Name:    "Push4",
		Summary: "Evaluates the next 4 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push4",
			Post: "A"},
		Parms:   []parm{dataParm{4}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x25,
		Name:    "Push5",
		Summary: "Evaluates the next 5 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push5",
			Post: "A"},
		Parms:   []parm{dataParm{5}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x26,
		Name:    "Push6",
		Summary: "Evaluates the next 6 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push6",
			Post: "A"},
		Parms:   []parm{dataParm{6}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x27,
		Name:    "Push7",
		Summary: "Evaluates the next 7 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push7",
			Post: "A"},
		Parms:   []parm{dataParm{7}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x28,
		Name:    "Push8",
		Summary: "Evaluates the next 8 bytes as a signed little-endian numeric value and pushes it onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "push8",
			Post: "A"},
		Parms:   []parm{dataParm{8}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x29,
		Name:    "PushB",
		Summary: "Pushes the specified number of following bytes onto the stack as a Bytes object.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "pushb 3 0x41 0x42 0x43",
			Post: `"ABC"`},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2A,
		Name:    "One",
		Synonym: "True",
		Summary: "Pushes 1 onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "one, true",
			Post: "1"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2B,
		Name:    "Neg1",
		Summary: "Pushes -1 onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "neg1",
			Post: "-1"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2C,
		Name:    "PushT",
		Summary: "Concatenates the next 8 bytes and pushes them onto the stack as a timestamp.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "pusht",
			Post: "timestamp A"},
		Parms:   []parm{dataParm{8}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2D,
		Name:    "Now",
		Summary: "Pushes the current timestamp onto the stack.",
		Doc:     "Note that 'current' may have special meaning depending on the context; in particular, repeated uses of this opcode may (and most likely will) return the same value within a given context.",
		Example: example{
			Pre:  "",
			Inst: "now",
			Post: "(current time as timestamp)"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2E,
		Name:    "PushA",
		Summary: "Evaluates a to make sure it is formatted as a valid ndau-style address; if so, pushes it onto the stack as a Bytes object. If not, error.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "pusha nda234...4b3",
			Post: "nda234...4b3"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x2F,
		Name:    "Rand",
		Summary: "Pushes a 64-bit random number onto the stack. Note that 'random' may have special meaning depending on context; in particular, repeated uses of this opcode may (and most likely will) return the same value within a given context.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "rand",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x30,
		Name:    "PushL",
		Summary: "Pushes an empty list onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "pushl",
			Post: "[]"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x40,
		Name:    "Add",
		Summary: "Adds the top two numeric values on the stack and puts their sum on top of the stack. attempting to add non-numeric values is an error.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "add",
			Post: "A+B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x41,
		Name:    "Sub",
		Summary: "Subtracts the top numeric value on the stack from the second and puts the difference on top of the stack. attempting to subtract non-numeric values is an error.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "sub",
			Post: "A-B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x42,
		Name:    "Mul",
		Summary: "Multiplies the top two numeric values on the stack and puts their product on top of the stack. attempting to multiply non-numeric values is an error.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "mul",
			Post: "A*B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x43,
		Name:    "Div",
		Summary: "Divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "div",
			Post: "int(A/B)"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x44,
		Name:    "Mod",
		Summary: "Divides the second numeric value on the stack by the top and puts the integer remainder on top of the stack. attempting to divide non-numeric values is an error, as is dividing by zero.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "mod",
			Post: "A % B"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x45,
		Name:    "DivMod",
		Summary: "Divides the second numeric value on the stack by the top and puts the integer quotient on top of the stack and the remainder in the second item on the stack. attempting to divide non-numeric values is an error, as is dividing by zero.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "divmod",
			Post: "A%B int(A/B)"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x46,
		Name:    "MulDiv",
		Summary: "Multiplies the third numeric item on the stack by the fraction created by dividing the second numeric item by the top; guaranteed not to overflow as long as the fraction is less than 1. An overflow is an error.",
		Doc:     "",
		Example: example{
			Pre:  "A B C",
			Inst: "muldiv",
			Post: "int(A*(B/C))"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x48,
		Name:    "Not",
		Summary: "If the top of the stack is 0, it is replaced by 1 -- otherwise, it is replaced by 0.",
		Doc:     "",
		Example: example{
			Pre:  "5 6 7",
			Inst: "not",
			Post: "5 6 0"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x49,
		Name:    "Neg",
		Summary: "The sign of the number on top of the stack is negated.",
		Doc:     "",
		Example: example{
			Pre:  "A",
			Inst: "neg",
			Post: "-A"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x4A,
		Name:    "Inc",
		Summary: "Adds 1 to the number on top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A",
			Inst: "inc",
			Post: "A+1"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x4B,
		Name:    "Dec",
		Summary: "Subtracts 1 from the number on top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A",
			Inst: "dec",
			Post: "A-1"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x4D,
		Name:    "Eq",
		Summary: "Compares (and discards) the two top stack elements. If they are equal in both type and value, leaves TRUE (1) on top of the stack, otherwise leaves FALSE (0) on top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "eq",
			Post: "FALSE"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x4E,
		Name:    "Gt",
		Summary: "Compares (and discards) the two top stack elements. If the types are different, fails execution. If the types are the same, compares the values, and leaves TRUE when: a) the top Number or Timestamp or ID is numerically greater than the second, b) the top list is longer than the second list, c) the top struct compares greater than the second by iterating in field order and using rules a) and b). .",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "gt",
			Post: "TRUE"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x4F,
		Name:    "Lt",
		Summary: "Like gt, using less than instead of greater than.",
		Doc:     "",
		Example: example{
			Pre:  "A B",
			Inst: "lt",
			Post: "FALSE"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x50,
		Name:    "Index",
		Summary: "Selects a zero-indexed element (the index is the top of the stack) from a list reference which is the second item on the stack (both are discarded) and leaves it on top of the stack. Error if index is out of bounds or a list is not on top of the stack.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z] 2",
			Inst: "index",
			Post: "Z"},
		Parms:   []parm{indexParm{"index"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x51,
		Name:    "Len",
		Summary: "Returns the length of a list.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z]",
			Inst: "len",
			Post: "3"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x52,
		Name:    "Append",
		Summary: "Creates a new list, appending the new value to it.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y] Z",
			Inst: "append",
			Post: "[X Y Z]"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x53,
		Name:    "Extend",
		Summary: "Generates a new list by concatenating two other lists.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y] [Z]",
			Inst: "extend",
			Post: "[X Y Z]"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x54,
		Name:    "Slice",
		Summary: "Expects a list and two indices on top of the stack. Creates a new list containing the designated subset of the elements in the original slice.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z] 1 3",
			Inst: "slice",
			Post: "[Y Z]"},
		Parms:   []parm{indexParm{"start"}, indexParm{"end"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x60,
		Name:    "Field",
		Summary: "Retrieves a field at index f from a struct; if the index is out of bounds, fails.",
		Doc:     "",
		Example: example{
			Pre:  "X",
			Inst: "field f",
			Post: "X.f"},
		Parms:   []parm{indexParm{"fieldindex"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x70,
		Name:    "FieldL",
		Summary: "Makes a new list by retrieving a given field from all of the structs in a list.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z]",
			Inst: "fieldl f",
			Post: "[X.f Y.f Z.f]"},
		Parms:   []parm{indexParm{"fieldindex"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x80,
		Name:    "Def",
		Summary: "Defines function block n, where n is a number larger than any previously defined function in this script. Function 0 is called by the system. Every function must be terminated by end, and function definitions may not be nested.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "def n",
			Post: ""},
		Parms:   []parm{functionIDParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x81,
		Name:    "Call",
		Summary: "Calls the function block, provided that its ID is greater than the index of the function block currently executing (recursion is not permitted). The function runs with a new stack which is initialized with the top n values of the current stack (which are copied, NOT popped). Upon return, the top value on the function's stack is pushed onto the caller's stack.",
		Doc:     "The function's return value is the top entry on its stack upon return.",
		Example: example{
			Pre:  "",
			Inst: "call n m",
			Post: ""},
		Parms:   []parm{functionIDParm{}, countParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x82,
		Name:    "Deco",
		Summary: "Decorates a list of structs (on top of the stack, which it pops) by applying the function block to each member of the struct, copying n stack entries to the function block's stack, then copying the struct itself; on return, the top value of the function block stack is appended to the list entry. The resulting new list is pushed onto the stack.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "deco n m",
			Post: ""},
		Parms:   []parm{functionIDParm{}, countParm{}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x88,
		Name:    "EndDef",
		Summary: "Ends a function definition; always required.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "enddef",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x89,
		Name:    "IfZ",
		Summary: "If the top stack item is zero, executes subsequent code. The top stack item is discarded.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "ifz",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x8A,
		Name:    "IfNZ",
		Summary: "If the top stack item is nonzero, executes subsequent code. The top stack item is discarded.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "ifnz",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x8E,
		Name:    "Else",
		Summary: "If the code immediately following an if was not executed, this code (up to end) will be; otherwise it will be skipped.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "else",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x8F,
		Name:    "EndIf",
		Summary: "Terminates a conditional block; if this opcode is missing for any block, the program is invalid.",
		Doc:     "",
		Example: example{
			Pre:  "",
			Inst: "endif",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x90,
		Name:    "Sum",
		Summary: "Given a list of numbers, sums all the values in the list.",
		Doc:     "",
		Example: example{
			Pre:  "[2 12 4]",
			Inst: "sum",
			Post: "18"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x91,
		Name:    "Avg",
		Summary: "Given a list of numbers, averages all the values in the list.",
		Doc:     "",
		Example: example{
			Pre:  "[2 12 4]",
			Inst: "avg",
			Post: "6"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x92,
		Name:    "Max",
		Summary: "Given a list of numbers, finds the maximum value.",
		Doc:     "",
		Example: example{
			Pre:  "[2 12 4]",
			Inst: "max",
			Post: "12"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x93,
		Name:    "Min",
		Summary: "Given a list of numbers, finds the minimum value.",
		Doc:     "",
		Example: example{
			Pre:  "[2 12 4]",
			Inst: "min",
			Post: "2"},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x94,
		Name:    "Choice",
		Summary: "Selects an item at random from a list and leaves it on the stack as a replacement for the list.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z]",
			Inst: "choice",
			Post: ""},
		Parms:   []parm{},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x95,
		Name:    "WChoice",
		Summary: "Selects an item from a list of structs weighted by the given field index.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z] f",
			Inst: "wchoice f",
			Post: ""},
		Parms:   []parm{indexParm{"fieldindex"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x96,
		Name:    "Sort",
		Summary: "Sorts a list of structs by a given field.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z] f",
			Inst: "sort f",
			Post: "The list sorted by field f"},
		Parms:   []parm{indexParm{"fieldindex"}},
		Enabled: true,
	},
	opcodeInfo{
		Value:   0x97,
		Name:    "Lookup",
		Summary: "Selects an item from a list of structs by applying the function block to each item in order, copying n stack entries to the function block's stack, then copying the struct itself; returns the index of the first item in the list where the result is a nonzero number; throws an error if no item returns a nonzero number.",
		Doc:     "",
		Example: example{
			Pre:  "[X Y Z]",
			Inst: "lookup n m",
			Post: "i"},
		Parms:   []parm{functionIDParm{}, countParm{}},
		Enabled: true,
	},
}
