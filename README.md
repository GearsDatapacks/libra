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

Integers can be of several types, varying in bit-size and signed-ness.
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
A `\` preceding a character that does not form an escape sequence will cause a syntax error.
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

### Unary Expressions
A unary expression consists of an operator and an operand. The operator is either prefix (comes before operand), or postfix (comes after).
Unary operators have only two precedence levels as they are controlled by order. Postfix operators have higher precedence than prefix ones.

Example:
```rust
mut i = 1
print(-i++) // -1
```

Here is a table of unary operators:

Operator|Function           |Position
--------|--------           |--------
++      |Increment          |Postfix
--      |Decrement          |Postfix
?       |Error propagation  |Postfix
!       |Error unwrapping   |Postfix
\+      |Arithmetic identity|Prefix
\-      |Arithmetic Negation|Prefix
!       |Logical Not        |Prefix
\*      |Pointer dereference |Prefix
&       |Reference          |Prefix
~       |Bitwise Not        |Prefix


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

## Pointers
A pointer is a value that holds an address of another value. This allows you to pass values around without copying them, or allowing functions to modify a value.
A pointer type is denoted: `*T`, where `T` is the type of the value it points to.
By default a pointer is immutable, meaning it cannot be used to modify the underlying value. A mutable pointer is denoted: `*mut T`. A mutable pointer is assignable to an immutable pointer.

You can create a pointer by referencing a value using the reference operator (`&`), and access the value it points to using the dereference operator (`*`).
A mutable pointer can be created by adding the `mut` keyword (`&mut <value>`). A mutable pointer can only reference a mutable value. 

Example:
```rust
fn increment(p: *mut i32) {
  *p += 1
}

mut i = 0
increment(&mut i)
print(i) // 1

// As opposed to
fn increment(mut i: i32) {
  i += 1 // i32 is passed by copy, so this doesn't modify the original value
}
mut i = 0
increment(i)
print(i) // still 0
```

The dereference operator doesn't always need to be used to access the value that a pointer references. Member access on a struct can still be used through a pointer.

Example:
```rust
struct Point {
  x, y: f32
}

fn (Point) sum(): f32 {
  return this.x + this.y
}

let point = Point {x: 3.2, y: 1.8}
let pointer_to_point: *Point = &point
let point_x = pointer_to_point.x // This works, no need to dereference
let sum = pointer_to_point.sum() // This also works
```

## Type-check Expressions
A type-check expression is used to determine if a value is of a certain type at runtime.
It uses the syntax: `<value> is <type>`.

Example:
```rust
union Int {
  i8, i16, i32, i64
}

fn get_bit_width(i: Int) -> i32 {
  if i is i8 {
    return 8
  } else if i is i16 {
    return 16
  } else if i is i32 {
    return 32
  } else {
    return 64
  }
}
```

## Cast expressions
Cast expressions can be used to convert a value from one type to another. This either changes the representation in memory or forces a union type to be of one type. The syntax is `<value> -> <type>`.

Example:
```rust
let int = 21
let float = int -> f32 // 21.0

union Thing {
  string, i32, f32
}

let my_thing: Thing = "Hello"
let my_str = my_thing -> string // Makes sure my_thing is a string
let my_int = my_thing -> i32 // ERROR: Thing is string not i32
```

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

If the function doesn't perform any intermediate operations, and simply returns a value, the `return` keyword can be omitted:
```rust
fn add(a, b: i32): i32 {a + b}
```

Functions can be called the function name, followed by `(`, `)` enclosing comma-separated arguments.
Function calls produce the value that the function returns, if any.

Example:
```rust
let result = add(1, 2) // result = 3
```

## Statements
Statements make up a Libra program, and don't produce a value.

### Blocks
A block statement is a grouping of statements enclosed in `{`, `}`. The statements are evaluated from top to bottom as in a regular program, but live in their own scope, meaning any variables created are not valid outside of it. This is useful when you want to perform an inline calculation without creating any unnecessary variables.

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

### If/else
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

### Return
A return statement is used to return a value from a function. The syntax is the `return` keyword followed by an optional expression.
The type of the expression must match the return type of the surrounding function (unless `void`, in which case there must not be a value), and the return statement must only be used within a function. Breaking either of these rules results in a compile-time error.

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
  print("Unknown task")
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

let value: module.Type = module.exported_function()
```

If you want to change the name of the imported module, you can use the `as` keyword.

Example:
```go
import "my_very_long_module_name" as mod
mod.foo()
```

You can specify a list of members to import directly into the current file using `import {...} from ...`, or import all exported member of that module using `*`.

Example:
```go
import {foo, bar} from "module"
foo()
bar()

import * from "module"
baz() // Exported by module
```

## Statements and expressions
Many statements can also be used as expressions. This allows them to produce values and be used wherever expressions would be.

### If/else expressions
Libra has no ternary operator. Instead, the if/else construct can be used as an expression. The resulting value can be obtained using the `yield` keyword.  
If there is only one expression in the block, the yield keyword can be omitted.  
If/else expressions don't require an else branch, but if one is not present, the resulting value is optional (`T?`).

Example:
```rust
let my_ternary = if foo {bar} else {baz}

let my_if_expr = if condition1 {
  do_thing()
  yield 1
} else if condition2 {
  do_other_thing()
  yield 2
} else {
  do_final_thing()
  yield 0
}
```

### Block expressions
Block expressions are just like if/else expressions, but they run a single block unconditionally. This can be used to calculate a value without polluting the scope with all the intermediate values used, and without creating a separate function for the logic.  

Example:
```rust
let result = {
  let x = 3.4
  let x_squared = x ** 2
  let two_x = x * 2
  let five = 5
  yield x_squared + two_x + five
}
// Only result is still in scope
```

### Function expressions
In Libra, functions are values. You can reference declared functions by name, but sometimes you want to create a function inline, without declaring is globally. For this you can use anonymous functions.  
Anonymous functions are just like regular functions except that they are created without a name and can be used inline as expressions.  
Additionally, anonymous functions can omit their return type as it can be inferred.

An example of how you might achieve some behaviour without anonymous functions:
```rust
fn my_filter(x: i32): bool {x % 2 == 0 || x < 5}

[1, 2, 4, 5, 6, 7, 91, 24].filter(my_filter)
// my_filter is now a global function
```

Here's how you could do it with anonymous functions:
```rust
[1, 2, 4, 5, 6, 7, 91, 24].filter(fn(x: i32) {x % 2 == 0 || x < 5})
// Anonymous function only exists within the scope of filter
```

A function value with arguments of type `A` and `B` and a return value of type `C` has the type: `fn(A, B): C`

### Loop expressions
In Libra, loops can produce values just like if statements. When loops are terminated using `break`, an optional value can be supplied. The loop expression will evaluate to this value.  
Since loops cannot be guaranteed to break with a value or even run at all, this returns an optional value. To remove the optional, use an `else` branch.

Example:
```rust
let result = for i in 1..10 {
  if i > 10 {
    break i
  }
} else {-1}
```

### Then/else
Then/else provides a neat way of handling optional values. Any expression that implements the `Chain` interface can be used with then/else.  
The `then` branch is run if the value is present, the `else` branch if otherwise. Examples of types that implement `Chain` by default are `T?` and `T!`.  
The `Chain` interface is defined as such:
```go
// P = Present, A = Absent
interface Chain<P, A> {
  chain(): P | A
}
```
Both branches are optional, but if the `else` branch is omitted, the resulting value will still be optional.  
Both branches can use block captures to receive the values from the optional value (`else` branch will receive a void value on `T?`).

Example:
```rust
fn safe_divide(a, b: f32): f32? {
  if b == 0 {
    return void
  }
  return a / b
}

fn unsafe_divide(a, b: f32): f32 {
  return safe_divide(a, b) else {
    panic("Tried to divide by zero")
  }
}

fn nan_divide_plus_one(a, b: f32): f32 {
  return safe_divide(a, b) then |result| {
    result + 1
  } else {
    f32.NaN
  }
}
```

Then/else can be used to implement the `?` and `!` operators:
```rust
fn foo(): i32? {
  return safe_divide() else {return void} + 1
}

fn bar(): i32 {
  return safe_divide() else {panic("Tried to unwrap void option")} + 1
}
```

Then/else can also be used with loops, to check whether the loop executed.

Example:
```rust
fn sum(values: i32[]): i32! {
  mut sum = 0
  for i in values {
    sum += i
  } then {
    return sum
  } else {
    return Err.new("Please provide at least one value to sum")
  }
}
```

## Type-declaration statements
Type-declaration statements are statements that can only appear at the top level of the program; that is, not within another statement such as a function.
They define a data-type that can be used in the program. Type names follow the PascalCase naming convention.

### Type aliases
A type alias simply binds a type to a name. It is useful for eliminating the need to repeatedly write a complex type.
Type aliases follow a similar syntax to variable declarations, and can also be referenced using their name.

Example:
```rust
type ComplexType = {string: i32[2]}[]
let complex_value: ComplexType = [{"values": [1,2]}]
```

### Structs
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

Structs can be defined without field names. This allows fields to be accessed the same a in tuples or arrays.

Example:
```rust
struct Coordinate3D { f32, f32, f32 }

let coords = Coordinate3D { 1.4, 82.3, -9.3 }
let x = coords[0] // 1.4
```

### Unit structs
A unit struct is a value with no size. This is equivalent to a union member with type void.

Example:
```rust
struct Unit

let empty_value = Unit
```

### Result and Option types
Results and options are ways of specifying types that might have some kind of error.
The option type is a union between `None`, a void value, and some type `T`. This is useful when a value might not be there, such as returning the element of a map. The syntax for an option type is `T?`.

A result type is a union between a type `T` and the `Error` tag. This allows for returning a value from a function that might error. The syntax for a result type is `T!`.

### Interfaces
An interface is a type that doesn't describe to a specific value, but rather a constraint on values.
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

fn (Robot) think(input: i32): i32 {
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

### Explicit interfaces
By default, any type with the methods described by an interface are assignable to that interface. A type only conforms to an explicit interface if all methods of that type required by the interface are tagged as implementing that interface.  
Using explicit interfaces is only recommended for simple interfaces which might pick up methods by chance due to common names.  

Example:
```rust
explicit interface Foo {
  foo(): i32
  bar(): f32
}

struct Fooer
fn (Fooer) foo(): i32 {10} // does not explicitly implement Foo
fn (Fooer) bar(): f32 {3.14}

struct Barer

@impl Foo
fn (Barer) foo(): i32 {Fooer.foo() - 3}
@impl Foo
fn (Barer) bar(): f32 {9.3}
```

#### impl blocks
If you need to conform to a large explicit interface, you will have to tag a lot of methods. To avoid this, you can use impl blocks. And impl block is a block that automatically tags all methods inside it.

Example:
```rust
explicit interface MyBigInterface {
  a(): A
  b(): B
  ...
  z(): Z
}

struct ThankGodForImplBlocks

@impl MyBigInterface {
  fn (ThankGodForImplBlocks) a(): A {...}
  fn (ThankGodForImplBlocks) b(): B {...}
  ...
  fn (ThankGodForImplBlocks) z(): Z {...}
}
```

### Tags
A tag is similar to an explicit interface, except that it has no required methods. Tags are implemented by anything which specifies it.

Example:
```rust
tag MyTag

@tag MyTag
struct MyTagImpl

let my_tagged_value: MyTag = MyTagImpl
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

// If we know that the type is Name, Name can be omitted from the value
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
```

Unions can also contain compound data-types, which can be defined within the body of the union.  
To access a type of the union, a cast can be used to try and cast to one of the specified members. Or, a shorthand syntax can be used.

```rust
union Property {
  Height: f32,
  Weight: f32,
  Age: u32,
  Name {
    first, last: string
  },
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

let my_prop = Property.Name {first: "John", last: "Doe"}
print((my_prop -> .Name).first) // John
// Shorthand syntax
print(my_prop.Name.last) // Doe
```

A shorthand for unions can be used: `<type1> | <type2> | ... | <type>`. This allows for inlining a union type without declaring a type for it, but it doesn't support naming of fields and therefore doesn't allow duplicate types.

Example:
```rust
mut int_or_float: i32 | f32 = 7
int_or_float = 12.3
```

### Untagged unions
Unions can also be marked as untagged, meaning the compiler doesn't store any information about which type is being represented. This is equivalent to C unions. Untagged unions cannot use the shorthand syntax.  
Untagged unions are unsafe and can cause invalid data such as dangling pointers.

Example:
```rust
@untagged
union IntPtr {
  int: usize,
  ptr: *i32
}

let my_fake_ptr = IntPtr.int(182)
let my_ptr = my_fake_ptr -> .ptr // my_ptr now points to address 182
// Or the shorthand:
let my_ptr2 = my_fake_ptr.ptr
let my_int: i32 = *my_ptr // Probably nonsense data
```

## Hello, world!
In Libra, no main function is needed. Any top-level statements will be run on program entry.
Therefore, Hello, world is just a single line of Libra:
```rust
print("Hello, world!")
```

## WIP
This section is a working in progress, and these concepts are all likely to change. This is simply somewhere for me to note down ideas as I come up with them; they will be refined later.

### Context
Every program has a context. This stores important information about the program state such as the currently active allocator.  
Context is passed by copy, to prevent called functions changing state that might mess up behaviour in the current function.

Example:
```rust
context.allocator = some_allocator
let alloced = alloc(i32)
// This is guaranteed not to change the current allocator, so it is safe to assume that context.allocator is still some_allocator
my_sub_fn()
free(alloced)
```

### Allocators
An allocator is a type that is responsible for allocating and freeing memory. Any type that stores data on the heap (such as lists or maps) uses an allocator to store that data.  
An allocator is any type that implements to the `Allocator` interface, defined in the `std:memory` library:
```go
interface Allocator {
  alloc(Type): *mut Type!
  free(*Type)
}
```

By default, all allocations in Libra use the default allocator, `CAllocator`. This uses `malloc` and `free` from LibC.  
You can change the allocator by modifying `allocator` field in the program context:
```rust
// Allocated using CAllocator
let my_4_ints: *mut i32[_] = alloc(i32[4])! // alloc is a builtin function that calls context.allocator.alloc
context.allocator = MyCustomAllocator{...}
// Allocated using MyCustomAllocator
let my_custom_4_ints: *mut i32[_] = alloc(i32[4])!
```

**Note**: The builtin `free` function calls the `free` method on the current allocator, so if the allocator has changed since something was allocated, the old allocator will need to be rememberer and used to free it.

Example:
```rust
let old_allocator = context.allocator
let old_alloced: *i32 = alloc(i32)!
context.allocator = other_allocator
let new_alloced: *i32 = alloc(i32)!
free(new_alloced) // This is fine, new_alloced was allocated using the current allocator
free(old_alloced) // This won't work since the current allocator didn't allocate old_alloced
// Do this instead:
old_allocator.free(old_alloced)
```

### Defer
A defer statement allows you to delay code execution until the end of a scope. This can be used, for example, to free anything allocated within a function.  
This is useful because it allows you to keep the code that frees data in the same place as where it is allocated, and also if you return from multiple places in the function, you only need to free it once.

For example, this code, which needs to free `my_thing` in three places:
```rust
let my_thing = alloc(Thing)

if foo {
  free(my_thing)
  return bar
} else if bar {
  free(my_thing)
  return foo
} else {
  free(my_thing)
  return baz
}
```

Can be replaced with this:
```rust
let my_thing = alloc(Thing)
// Gets automatically freed in all three cases
defer free(my_thing)

if foo {
  return bar
} else if bar {
  return foo
} else {
  return baz
}
```

### Type parameters
Functions can take type parameters, parameters of type `const Type`, which can be used in the signature of the function.

Example:
```rust
fn add(T: const Type, a, b: T): T {
  if T == i32 {
    return add_i32(a, b)
  } else if T == f32 {
    return add_f32(a, b)
  } else {
    panic("Cannot add non-number types")
  }
}
```

In many cases, type parameters can be inferred, and the shorthand syntax can be used: `$T`.

Example:
```rust
fn add(a, b: $T): T {
  // Body is the same
  ...
}

add(1, 2) // T can be inferred as i32
add(3.1, 2.5) // T can be inferred as f32
```

Custom types can also take type parameters. Since type parameters of custom types can't always be inferred, there is no shorthand syntax, but the parameters are optional if it they can be.

Example:
```rust
union Option(T) {
  T, void
}

let my_option: Option = 1 // T can be inferred as i32
let my_other_option: Option(string) = void // T cannot be inferred and must be provided
```

Type constraints can be added to type parameters to limit what types can be given.

Example:
```rust
fn add(a, b: $T: Add): T {
  return a + b
}
```

```rust
import {Add, Subtract as Sub, Multiple as Mul, Divide as Div} from "math"

enum Operation {Add, Sub, Mul, Div}
struct Expression(T: Add & Sub & Mul & Div) {
  left, right: T,
  operation: Operation
}

fn (Expression($T)) calculate(): T {
  return switch this.operation {
    case .Add => this.left + this.right
    case .Sub => this.left - this.right
    case .Mul => this.left * this.right
    case .Div => this.left / this.right
  }
}
```

### Builder pattern
Libra supports the builder pattern at a language level, allowing structs to provide builder methods instead of being constructed in the usual way.

```rust
struct Averages {
  ~values: i32[]
  mean: f32,
  mode, median: i32,
}

fn Averages.new(values: i32[]): ~Averages { Averages { values } }

fn (~Averages) values(values: i32[]): ~Averages {
  this.values.extend(values)
  return this
}

fn (~Averages) value(value: i32): ~Averages {
  this.values.push(value)
  return this
}

fn (~Averages) calculate(): Averages {
  mut frequencies: {i32: i32} = {}
  mut sum = 0

  for value in values {
    frequencies[value]++
    sum += value
  }

  this.mode = 0
  mut max_frequency = 0
  for (value, freq) in frequencies {
    if freq > max_frequency {
      this.mode = value
      max_frequency = freq
    }
  }

  let len = this.values.len()
  this.median = if len % 2 == 0 {
    (this.values[len / 2] + this.values[len / 2 + 1]) / 2
  } else {
    this.values[len / 2 + 0.5]
  }

  this.mean = sum / count
  return this
}

let averages = Averages.new([1, 2, 3]).values([4, 5, 6]).value(7).calculate()
let mean = averages.mean
// ERROR
averages.value(10)
```
