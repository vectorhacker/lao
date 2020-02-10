Lao interpreter
===============

This is an interpreter implemented in Go for a simple BASIC-like language called Lao. This was created as part of a programming languages course.

Parts
-----

This is implemented in the following simple parts:

1. Tokenizer - turns stream of characters into recognizable tokens
2. Parser - Parser turns the tokens into an AST
3. Interpreter - Takes in an AST and runs the program.


How to use
----------

In order to install and run a sample program:

```bash

go get github.com/vectorhacker/lao
go install github.com/vectorhacker/lao

lao <path_to_program>
```

Next steps
----------

~~The Lao language is incomplete as a language. It is missing simple jumps and loops.~~ It also lacks a repl.

To be implemented
---------------

1. Repl
2. ~~Goto statement~~
3. ~~Label statement~~
