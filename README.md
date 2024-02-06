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
	// Mathematical expressions to be evaluated
	expression := "(2 + 3) * 4 ^ 2 / (5 - 1)"
	expression2 := fmt.Sprintf("(%d + %f) * 4 ^ 2 / (%f - %f)", 2, 3.5, 5.7, 1.0)

	// Parse and evaluate the expression
	result, err := parser.Parse(expression)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// Print the result: 20.000000
	fmt.Printf("Result of '%s': %f\n", expression, result)

	// Parse and evaluate the expression
	result2, err := parser.Parse(expression2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result: 18.723404255
	fmt.Printf("Result of '%s': %.9f\n", expression2, result2)
}