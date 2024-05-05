# Libra
A programming language which balances low-level control and easy abstraction.

## Syntax

### Statements and Expressions
In Libra, there are expressions and statements. Statements are the components in a program, containing expressions within them. Most expressions produce values.
Expressions can appear on their own wherever a statement is allowed.

Most statements can appear anywhere in a program, except for a few which must appear at the top-level of a file, such as imports and type declarations.

### Comments
Any text following `//` in libra is ignored, until the end of a line.
Multiline comments can also be defined using `/*` and `*/`

### Whitespace
All whitespace in Libra is ignored, except for newlines, with separate statements.

### Semicolons and Newlines
In Libra, semicolons are not required, and statements can simply be separated by newlines. However, semicolons can also be used in place of a newline.

Semicolons behave almost identically to newlines, except that they completely cut off any statement or expression; whereas newlines allow, where unambiguous, the continuation of expressions and statements on the next line.

Example:
```rust
1
+ 2 // This is ambiguous, and treated as a different expression

1 +
2 // This is unambiguous, and is treated as one expression

1 +; 2 // This semicolon forces the expression to end, causing a syntax error
```

### Trailing commas
Whenever an expression or statement requires something separated by commas, a comma is always allowed after the last item in the list.

Example:
```rust
[1, 2, 3,]
```

## Types and Values
In Libra, most expressions produce values. A value is a piece of data stored in memory, that represents something.

Every value has a type, which dictates how the value is interpreted, and what operations can be performed on it.

There are a set of primitive types built in to the Libra programming language, as well as whatever types the program has defined.

For example, the value `42` is of type `i32`, a primitive Libra type which represents a 32-bit signed integer.

## Literal Expressions
A literal expression encodes its value in itself, its value is known at compile-time, and it has a single primitive type.

### Booleans
A boolean is either `true` or `false`, representing those values.
Booleans are of type `bool`.

### Integers
An integer is a number without a decimal point. 
Integers can be separated with `_` to improve their readability.
Integers can also be prefixed with `0b`, `0o` or `0x` to be encoded in base-2, base-8, or base-16 respectively.

Example:
```rust
1927
1_000_000
0b1001
0o1702
0xc0de10
```

Integers must not end with the delimiter (`_`). For example `1_000_` produces a syntax error.

Integers can be of several types, varying in bit-size and signded-ness.
For example, the `u16` type represents a 16-bit unsigned integer,
or the `i8` type, which represents a 8-bit signed integer.
Valid bit-sizes are powers of 2 between 8 and 64 inclusive.


### Floating point numbers
A floating point number (float) is one with a decimal point.

Floats can also be separated by `_`, but must not end with them.
Floats must have at least one digit on each side of the decimal point.
Floating point values can only be expressed in base-10.

Example of valid floats:
```rust
17.3
0.61
2.0
3.141_592_653
1_2.3_4
```

Example of invalid floats:
```rust
.92
12.
17_.3
82._15
1_2.3_4_
0b10.1001
```

Floats can be of varying bit-width, just like integers.
For example, the `f32` type means a 32-bit floating point value.

### Strings
A string represents a piece of text.
Strings are delimited with `"`.
Strings are of type `string`.

Strings can span multiple lines.
If you want to use the `"` character in a string, you can escape it with a `\`.
If you want to use a `\` character, it must be double escaped: `\\`.
The `\` character can also be usd in combination with other characters to create various escape sequences, e.g. `\n` (line feed) or `\t` (horizontal tab).
A `\` preceeding a character that does not form an escape sequence will cause a syntax error.
A string with no terminating `"` will also cause a syntax error

Example:
```rust
"This is a plain string"
"This string contains a \""
"This string does not escape the terminating quote \\"
"This string contains\na newline "
"This string
also contains a newline"
```

## Compound expressions
A compound expression is an expression containing one or more sub-expressions. Its value is known at compile-time if all sub-expressions are known at compile-time.
A compound expression's type depends on the types of its sub-expressions.

### Arrays and Lists
An array is an ordered list of values of one type.
It is defined using `[` and `]` containing comma-separated expressions.

Example:
```rust
[]
[1, 2, 3]
[true, true, false, true]
["hi", "there"]
```

There are two types that arrays can be: arrays and lists.

An array has a fixed length which must be known at compile-time.
A list has a variable length and can grow and shrink throughout the program.

An array containing elements of type `T` has the type: `T[<len>]`, where the value of `<len>` must be known at compile-time. If the length can be inferred, you can use `_` as the length.

A list containing elements of type `T` has the type `T[]`.

### Maps
A map is a a hashmap of key-value pairs. The keys must be of one type, and values of one type.
A map is expressed with `{` and `}`, with `key: value` pairs separated by commas.

Example:
```javascript
{}
{"a": 1, "b": 2, "z": 26}
{
  1: true,
  3: true,
  8: false,
}
{ true: "true", false: "false" }
```

A map with keys of type `K` and values of type `V` has the type: `{K: V}`


### Tuples
A tuple is a set of values, fixed in length, of possibly different types. The types of each value in the set must be known at compile-time.
Syntactically, a tuple is a set of comma-separated expressions, surrounded by `(`, `)`

Example:
```rust
(1, true, "Hello")
(1, ("2", 3))
```

A tuple containing values of type `A`, `B` and `C` has the type: `(A, B, C)`.

### Ranges
A range represents a range of values, such as integers. It has a start and an end, and is expressed using the `..` operator.

Example:
```rust
1..10
-12..33
1+2..(7-3)
```

### Binary Expressions
A binary expression is an expression with an operator and two operands. Each operator specifies what operation to perform on the operands, and produces the resulting value.

The resulting type of a binary expression depends on the operator and operands.

Example:
```rust
1 + 2
true && false
```

Each operator has a precedence value, and operators of higher precedence bind more tightly to their operands.
If two operators have equal precedence, the expression is evaluated according to that operator's associativity (If left, left-to-right, if right, right-to-left).

Example:
```rust
1 + 2 * 3 // * has a higher precedence than`+, so this is evaluated: 1 + (2 * 3)
1 + 2 - 3 // + has equal precedence to -, so it is evaluated left-to-right
```

Here is a table of operators in Libra:

Operator|Precedence|Function      |Associativity
--------|----------|--------------|-------------
=       |1         |Assignment    |right
+=      |1         |Add-Assign    |right
-=      |1         |Sub-Assign    |right
\*=     |1         |Mul-Assign    |right
/=      |1         |Div-Assign    |right
%=      |1         |Mod-Assign    |right
&&      |2         |Logical And   |left
\|\|    |2         |Logical Or    |left
<       |3         |Less Than     |left
<=      |3         |LT Eq         |left
\>      |3         |Greater Than  |left
\>=     |3         |GT Eq         |left
==      |3         |Equality      |left
!=      |3         |Inequality    |left
\.\.    |4         |Range         |left
<<      |5         |Left-Shift    |left
\>\>    |5         |Right-Shift   |left
&       |5         |Bitwise And   |left
\|      |5         |Bitwise Or    |left
\+      |6         |Addition      |left
\-      |6         |Subtraction   |left
\*      |7         |Multiplication|left
/       |7         |Division      |left
%       |7         |Modulo        |left
**      |8         |Exponentiation|right

You can also use `(`, `)` to group an expression, overriding precedence.

Example:
```rust
1 + 2 * 3 // Evaluated as 1 + (2 * 3) = 7
(1 + 2) * 3 // Evaluated a (1 + 2) * 3 = 9
```

## Variables
A variable is a named reference to a value.
A variable's name must begin with a letter or `_`, and continue with letters, numbers and `_`.
Variables in Libra follow the `snake_case` convention.

Variables can be declared in three different ways (all of which are statements):
```rust
// The value of a must be known at compile-time and cannot change
const a = 1

// The value of b cannot change (immutable)
let b = 2

// The value of c can change (mutable)
mut c = 3
```

In any case, the value referenced by a variable can be used as an expression, simply by the name of that variable:

```rust
let a = 1
let b = a + 1 // Evaluates to 2
```

If a variable is mutable, it can be re-assigned, but the value must be of the same type.

Example:
```rust
mut foo = 1
foo = 2 // Foo is now 2

// Equivalent to foo = foo + 1
foo += 1 // Foo is now 3
foo = 3.9 // Error: trying to assign float to int
```

When declaring variables, you can specify the type of the variable, either if it is not the default type of the value given, or for ease of reading.

Example:

```rust
// Type is provided here for readability, it changes nothing
const num: i32 = 42

// Forces value to be unsigned, would be signed by default
let unsigned: u32 = 20

// Forces value to be smaller bit-width, would be 32 by default
let small: i8 = 5

// Stores 1 as a float rather than int
mut float: f64 = 1

// Changes array to list, would be fixed size by default
mut arr: f32[] = [1.3, 3.14, -6.9]
```

If the specified type for a variable does not match the type of the value assigned to it, it produces an error.

## Functions
A function is a piece of code that can be called from anywhere, allowing behaviour to be determined by the arguments to be passed to it.

Functions are defined using the `fn` keyword, followed by the name of the function. This is a function declaration statement.
Then comes the arguments, a comma-separated list of names, with their types. If the type of an argument is omitted, it assumes the type of the following argument. This means that the last argument of a function must have a type annotation.
After the arguments, comes the return type. This can be omitted if the function does not return a value.
Finally is the body of the function, the code that is run when called. All arguments can be used as variable in this block.

```rust
fn add(a, b: i32): i32 {
  let c = a + b
  return c
}
```

Functions can be called the function name, followed by `(`, `)` enclosing comma-separated arguments.
Function calls produce the value that the function returns, if any.

Example:
```rust
let result = add(1, 2) // result = 3
```

// WhileLoop
// ForLoop
// FunctionDeclaration
// ReturnStatement
// TypeDeclaration
// StructDeclaration
// InterfaceDeclaration
// ImportStatement
// EnumDeclaration

## Statements
Statements make up a Libra program, and don't produce a value.

### Block statements
A block statement is a grouping of statements enlosed in `{`, `}`. The statements are evaluated from top to bottom as in a regular program, but live in their own scope, meaning any variables created are not valid outside of it. This is useful when you want to perform an inline calculation without creating any unnecessary variables.

Example:
```rust
mut value = 1
{
  const two = 2
  let other_thing = two + 73
  value += other_thing
}
print(value) // 76
print(other_thing) // ERROR: other_thing is now out of scope
```

### If/else statements
An if statement is used to run a piece of code based on a condition. It takes a condition expression (which must be of type `bool`) and a block statement to run as the body.
If the condition succeeds, it runs the block. Otherwise, if there is an else branch, it runs that. An else branch is the `else` keyword followed by another conditional if statement, or an unconditional block statement.

Example:
```rust
if thing_is_true() {
  print("true")
} else {
  print("false")
}

if a > 0 {
  print("a is positive")
} else if a < 0 {
  print("a is negative")
} else {
  print("a is zero")
}
```

### While loops
A while loop is used to repeatedly execute a block of code while a condition is true. The syntax is identical to an if statement, except using the `while` keyword.

Example:
```rust
mut i = 1

while i < 10 {
  print(i) // Prints 1 (inclusive) to 10 (exclusive)
  i++
}
print("Done!")
```

### For loops
A for loop iterates over values in an iterator. It is used to, for example, perform an action for each item in an array or other sequence.
A for loop is defined using the `for` keyword, followed by the name of the variable that holds the element of the iterator each iteration, then the `in` keyword, the value to iterate over, and a block to run each iteration.

Example:
```rust
// Runs for every number between 1 (inclusive) and 10 (exclusive)
for i in 1..10 {
  print(i)
}

const primes = [2, 3, 5, 7, 11, 13]
for prime in primes {
  print(prime, "is prime")
}
```

### Break and Continue
Sometimes you want to exit a loop before it normally should. You can use the `break` keyword for that.
Or, if you just want to skip to the next iteration of the loop, use `continue`.
Both break and continue are disallowed outside of loops, and produce a compile-time error if used that way.

Example:
```rust
for i in 1..10 {
  if i == 5 {
    break
  }
  print(i) // prints 1 to 4 (inclusive)
}

for i in 1..10 {
  if i % 2 == 0 {
    continue
  }
  print(i) // prints 1, 3, 5, 7, 9
}
```

### Return statements
A return statement is used to return a value from a function. The syntax is the `return` keyword followed by an optional expression.
The type of the expression must match the return type of the surrouding function (unless `void`, in which case there must not be a value), and the return statement must only be used within a function. Breaking either of these rules results in a compile-time error.

Example:
```rust
fn add(a, b: i32): i32 {
  return a + b
  print("This statement isn't reached")
}
print(add(1, 2)) // 3

fn perform_task(task: string) {
  if task == "greet" {
    print("Hello!")
    return
  }
  print("Unkown task")
}
perform_task("greet")
```

### Imports
When writing Libra code, you will often want to split code into multiple files/modules. An import can be used to bring in code from other modules.
A basic import statement is simply the `import` keyword followed by the path to the module to import.
This defines a struct containing all exported members from that module, and names it after the module name.

Example:
```go
import "path/to/module" // Gets imported as "module"

let value = module.exported_function()
```

If you want to change the name of the imported module, you can use the `as` keyword.

Example:
```go
import "my_very_long_module_name" as mod
mod.foo()
```

You can specify a list of members to import directly into the current file using `import {...} from ...`, or import all exported memeber of that module using `*`.

Example:
```go
import {foo, bar} from "module"
foo()
bar()

import * from "module"
baz() // Exported by module
```

## Type-declaration statements
Type-declaration statements are statements that can only appear at the top level of the program; that is, not within another statement such as a function.
They define a data-type that can be used in the program. Type names follow the PascalCase naming convention.

### Type aliases
A type alias simply binds a type to a name. It is useful for eliminating the need to repeatedly write a complex type.
Type aliases follow a similar syntax to variable declarations, and can also be referenced using their name.

Example:
```rust
type ComplexType = {string: [2]i32}[]
let complex_value: ComplexType = [{"values": [1,2]}]
```

### Struct declarations
A struct declaration creates a custom type, which is a compound of named values.

Example:
```rust
struct Person {
  name: string,
  age: u8
}
```

This defines a struct `Person`, with the fields name (a string), and age (an unsigned 8-bit integer).
`Person` is now a data-type that can be instantiated, and used as a value.

Example:
```rust
// Instantiate the Person struct
let bob = Person {
  name: "Bob",
  age: 42
}

// Access the name field of bob
let bob_name = bob.name // "Bob"
```

### Interface declarations
An interface is a type that doesn't describe to a specifc value, but rather a constraint on values.
An interface simply defines a set of methods a type must have to conform to it. Any type with those methods automatically is assignable to that interface type.

Example:
```rust
// Anything with these two methods is a valid Sentient
interface Sentient {
  think(i32): i32,
  feel(): string
}

struct Human {
  name: string,
  emotion: string,
}

fn (Human) think(input: i32): i32 {
  return input + 1
}
fn (Human) feel(): string {
  return this.emotion
}

struct Robot {
  state: i32
}

fn (Robot) thing(input: i32): i32 {
  this.state += input
  return this.state
}
fn (Robot) feel(): string {
  return "Beep boop. Robots can't feel"
}

fn live(being: Sentient) {
  print(being.think(7))
  print(being.feel())
}

// Valid! Person implements Sentient
live(Person { name: "Julie", emotion: "Happy" })
// Valid! Robot implements Sentient
live(Robot { state: 94 })
// Invalid! i32 doesn't have either method required by Sentient
live(7)
```

### Enums
An enum is a type that is restricted to a set of possible values.
An enum has an underlying type, and each variant has a value of that type. By default, the underlying type is `i32`; if the type is a kind of integer, it automatically assigns incrementing values to each of the variants. If the type is a string, it automatically has the value of the name of that member. Otherwise, the values must be assigned manually.

Example:
```rust
enum Colour {
  Red, // 0
  Green, // 1
  Blue, // 2
  Pink = 105, // 105
  Yellow, // 106
}

let my_colour = Colour.Pink
if my_colour == Colour.Blue {
  print("My favourite colour!")
} else {
  print("The colour is adequate")
}

enum Name: string {
  Bob = "Bob",
  Ann = "Anne",
  Richard, // "Richard"
}

// If we know that the type is Name, Name can be ommited from the value
let person_name: Name = .Richard
```

You can use the static `from` method automatically generated for enums to construct a value for that enum from a raw value of that type, but since the values are limited it returns an optional result.
If you want to get the underlying value of an enum member, you can use the `raw` method.

Example:
```rust
let name = Name.from("Jane")
if name.some() {
  print("Jane is a valid name")
} else {
  print("Please choose a different name")
}

print(Colour.Blue.raw()) // 2
```

### Unions
A union is a value that can be of multiple types. It is tagged at runtime so the language knows which type it is.
If a union only has one member of a type, it can be inferred, otherwise an explicit member needs to be specified.

Example:
```rust
union Number {
  u8, u16, u32, u64,
  i8, i16, i32, i64,
  f8, f16, f32, f64,
}

// Number.f32
mut num: Number = 15.6
// Number.i32
num = 7

union Property {
  Height(f32),
  Weight(f32),
  Age(u32)
}

// Can be inferred: only one member has type u32
let age: Property = 102
// Can't be inferred: two members have type f32, so an explicit member must be specified
// But, since we explicitly use Property, the type annotation can be omitted
let height = Property.Height(1.67)
// If the type is known, it can be omitted from the value
let weight: Property = .Weight(61.3)

print(height == .Weight(1.67)) // False: one is weight, one is height
print(age == .Age(102)) // True: they completely match
```

## Hello, world!
In Libra, no main function is needed. Any top-level statements will be run on program entry.
Therefore, Hello, world is just a single line of Libra:
```rust
print("Hello, world!")
```
