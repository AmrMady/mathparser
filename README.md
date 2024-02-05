# mathparser

`mathparser` is a simple Go package that provides a mathematical expression parser. It can parse and evaluate mathematical expressions with basic operators and functions.

## Installation

To use this package in your Go project, you can install it using the following `go get` command:

```bash
go get github.com/AmrMady/mathparser/parser
```

## Example usuage
```go
package main

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser"
)

func main() {
	// Mathematical expression to be evaluated
	expression := "(2 + 3) * 4 ^ 2 / (5 - 1)"

	// Parse and evaluate the expression
	result, err := parser.Parse(expression)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result
	fmt.Printf("Result of '%s': %f\n", expression, result)
}