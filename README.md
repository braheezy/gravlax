# gravlax

A Go implementation of an interpreter for the educational Lox programming language as taught in the book [Crafting Interpreters](https://craftinginterpreters.com/).

The second half of the book covers a bytecode VM approach to implementing Lox. I did that in Zig [here](https://github.com/braheezy/zig-lox).

## Usage
You need Go 1.21.9+.

Clone the code and get inside the directory:

```bash
git clone https://github.com/braheezy/gravlax
cd gravlax
```

Run the REPL (Ctrl+D to exit):

```bash
$ go run main.go
> print("hey there");
hey there
> (ctrl+d)
bye!
```
Or run a file:

```bash
$ echo 'print("hey there");' > hello.lox
$ go run main.go hello.lox
hey there
```
## Notes
The tutorial book covers `jlox`, a Java implementation using a tree-walk interpreter approach to executing Lox programs. `gravlax` is the same thing, but in Go.

Of note, this implementation:
- supports block comments
- the `break` keyword
