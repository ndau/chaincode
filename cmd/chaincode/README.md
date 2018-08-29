# Chaincode: a programming language for expressions

The purpose of Chaincode is to allow a limited scripting language for expressions. The key element is that it is explicitly not Turing-complete, which in practical terms means that it has no mechanisms for looping or for constructs that allow nonterminating programs to be created. It is amenable to static analysis.

Basically, it is a language and virtual machine designed for creating expressions. It has conditionals and functions, but only in very limited form.

## Basic syntax

* A source file is composed of a some constant definitions (some constants are predefined), some handlers, and some functions
* Functions look like `function name(arg:type, arg2:type2) : rettype { ... }
* Functions must return one value (that value may be a list)
* The parameter list and return value of the main() function is defined in the context
* There are some predefined functions (which map to opcodes)
* Conditionals require parentheses and must evaluate to a boolean (there are no "truthy" or "falsy" values)
* Conditionals and functions require curly braces
* Every statement must end with a semicolon
* Assignments are not expressions
* Variables are scoped at the function level -- each function has only its own parameters and variables to play with
* Constants are defined at the global level only
* Constants are textual substitutes for values

## The requirements

* structs must be defined
* functions take typed parameters
* functions must be defined in the appropriate order (only forward function calls are permitted)
* main() must be the first function and its parameters are defined by the context
* Runtime errors cause the script to terminate

* Execution starts at main(), which must be the first function in the file, and its parameters must match those of the defined context.
*

Account scripts require an extra byte of context, which is used to indicate the minimum number of keys required for this script to be valid.

