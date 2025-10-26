# SILD - A TypeScript to Go transpiler

SILD is a TypeScript to Go transpiler that is written in Go. Simple enough.

## Why?

We're living in an era of countless JavaScript runtimes - blazing fast, but not
native. Go offers both speed and native compilation. SILD bridges TypeScript's
type safety with Go's performance characteristics. So while I was learning Go,
I decided to create a transpiler that could convert TypeScript to Go.

## Architecture

SILD uses a three-phase compilation process:

1. **Scanner**: Tokenizes TypeScript source code
2. **Parser**: Builds an Abstract Syntax Tree using Recursive Descent parsing
3. **Code Generator**: Emits idiomatic Go code from the AST

## Installation & Usage

Build:

```bash
git clone https://github.com/toyaAoi/sild.git
cd sild
make
```

Install:

```bash
sudo make install
```

Uninstall:

```bash
sudo make uninstall
```

Usage:

```bash
sild -o <output_file> <input_file>
```

## Examples

### Simple Number Assignment

```typescript
let x: number = 42;
```

```bash
sild -o main.go main.ts
```

```go
package main

func main() {
    x := 42
}
```

### Function Declaration

```typescript
function add(a: number, b: number): number {
  return a + b;
}

let result: number = add(1, 2);
```

```bash
sild -o main.go main.ts
```

```go
package main

func add(a int, b int) int {
    return (a + b)
}

func main() {
    result := add(1, 2)
}
```

### Multiple Functions and Variables

```typescript
function add(a: number, b: number): number {
  return a + b;
}

function multiply(a: number, b: number): number {
  return a * b;
}

let x: number = 5;
let y: number = 10;

let resultAdd: number = add(x, y);
let resultMultiply: number = multiply(x, y);
```

```bash
sild -o main.go main.ts
```

```go
package main

func add(a int, b int) int {
    return (a + b)
}

func multiply(a int, b int) int {
    return (a * b)
}

func main() {
    x := 5
    y := 10
    resultAdd := add(x, y)
    resultMultiply := multiply(x, y)
}
```

## Limitations

- Variables should be declared with `let` and explicitly typed
- Type inference is not supported
- Only supports basic types: number, string, boolean
- Only supports basic operators: +, -, \*, /, %
- Error handling needs improvement

## Roadmap

- [x] Variables
- [x] Functions
- [x] Expressions
- [ ] Control Flow
- [ ] Loops
- [ ] Conditionals
- [ ] Arrays
- [ ] Objects
