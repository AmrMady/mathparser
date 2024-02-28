# mathparser

`mathparser` is a simple Go package that provides a mathematical expression parser. It can parse and evaluate mathematical expressions with basic operators and functions.

## Installation

To use this package in your Go project, you can install it using the following `go get` command:

```bash
go get github.com/AmrMady/mathparser/parser
```

## Example 1
```go
package main

import (
"fmt"
"github.com/AmrMady/mathparser/parser"
)

func main() {
	// Simple mathematical expression to be evaluated
	expression := "(2 + 3) * 4 ^ 2 / (5 - 1)"

	// Parse and evaluate the expression using ParseSimple
	result, err := parser.ParseSimple(expression)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result: 20.000000
	fmt.Printf("Result of '%s': %f\n", expression, result)


	// Simple mathematical expression to be evaluated
	expression2 := fmt.Sprintf("(%d + %f) * 4 ^ 2 / (%f - %f)", 2, 3.5, 5.7, 1.0)
	
	// Parse and evaluate the expression
	result2, err := parser.ParseSimple(expression2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result: 18.723404255
	fmt.Printf("Result of '%s': %.9f\n", expression2, result2)
}

```


## Example 2
```go
package main

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser"
)

func main() {
	// Mathematical expression with variables to be evaluated
	expression := "x * y + (z - w) / a"

	// Define variables and their values
	variables := map[string]float64{
		"x": 5.5,
		"y": 4,
		"z": 20,
		"w": 15,
		"a": 2.5,
	}

	// Parse and evaluate the expression with variables
	result, err := parser.ParseWithVariables(expression, variables)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the result: 110
	fmt.Printf("Result of '%s': %f\n", expression, result)
}
```