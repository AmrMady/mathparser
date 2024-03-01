package advancedMathematicalFunctions

import (
	"fmt"
	"math/big"
)

// IsPrime checks if a number is prime
func IsPrime(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("isPrime requires exactly one argument")
	}

	// Convert *big.Float to *big.Int for IsPrime method
	n, _ := args[0].Int(nil)
	result := big.NewInt(0)
	if n.ProbablyPrime(0) { // Using ProbablyPrime with 0 for simplicity; adjust as needed
		result.SetInt64(1) // Return 1 for prime
	}
	return new(big.Float).SetInt(result).SetPrec(precision), nil
}

// GCD / Greatest Common Divisor calculates the greatest common divisor of two numbers
func GCD(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("gcd requires exactly two arguments")
	}

	a, _ := args[0].Int(nil)
	b, _ := args[1].Int(nil)
	gcd := new(big.Int).GCD(nil, nil, a, b)
	return new(big.Float).SetInt(gcd).SetPrec(precision), nil
}
