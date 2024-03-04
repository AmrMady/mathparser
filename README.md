
# mathparser 🧮✨

`mathparser` is an advanced Go package designed to parse and evaluate mathematical expressions with flair! It supports basic arithmetic operations, variable substitution, and returns results in various formats including float64, string, *big.Float, and byte array without losing precision.

## Installation 💾

To include `mathparser` in your Go project, install it with the `go get` command:

```bash
go get github.com/AmrMady/mathparser/parser
```

## Usage 📘

This package offers multiple functions to parse and evaluate mathematical expressions:

### ParseSimple

Evaluates simple expressions and returns the result as a float64.

```go
result, err := parser.ParseSimple("2 + 3 * 4")
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println("Result:", result)
// Output: Result: 14
```

### ParseWithVariables

Evaluates expressions with variables, returning the result as a float64.

```go
variables := map[string]float64{"x": 4, "y": 5}
result, err := parser.ParseWithVariables("x * y + 2", variables)
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println("Result:", result)
// Output: Result: 22
```

### ParseAndReturnString

Returns the result as a string to ensure precision.

```go
result, err := parser.ParseAndReturnString("1 / 3")
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println("Result:", result)
// Output: Result: 0.3333333333333333
```

### ParseAndReturnBigDecimal

Returns the result as a *big.Float for high precision calculations.

```go
result, err := parser.ParseAndReturnBigDecimal("1 / 3")
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println("Result:", result)
// Output: Result: 0.3333333333333333...
```

### ParseAndReturnBytes

Returns the result as a byte array, useful for serialization.

```go
bytesResult, err := parser.ParseAndReturnBytes("1 / 3")
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println("Bytes result:", bytesResult)
// Output: Bytes result: [byte array representing the big.Float]
```

## Complex Expressions Examples

#### Complex Algebraic Expression

```go
expression := "(12^3 - 2^7) / (4.5 + 3.5^2 - sqrt(2))"
// Expected Output: 104.33113466384367
```

#### Trigonometric Expression

```go
expression := "sin(45) + log(100) * tan(30)"
// Using degrees. Convert to radians if necessary in your implementation.
// Expected Output: 3.365903027730811
```

#### Leibniz Formula for Pi (Simplified)

```go
expression := "4 * sum((-1)^n / (2*n + 1) for n in range(1001))"
// Simplified approximation of Pi.
// Expected Output: 3.1425916543395442
```

These examples showcase the package's flexibility and capability in handling complex mathematical expressions, including an iconic equation related to Pi.

## Contributing

Contributions to `mathparser` are welcome! Please submit a pull request or open an issue to discuss proposed changes or additions.
