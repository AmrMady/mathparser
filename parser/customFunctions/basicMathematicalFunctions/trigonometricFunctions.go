package basicMathematicalFunctions

import (
	"fmt"
	"math"
	"math/big"
)

func Sin(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sin requires exactly one argument")
	}
	x := args[0]
	result := new(big.Float).SetPrec(precision)
	xPowerI := new(big.Float).SetPrec(precision).Copy(x) // x^i, starts as x^1
	factorialI := big.NewInt(1)                          // i!, starts as 1!

	for i := int64(1); ; i += 2 {
		term := new(big.Float).SetPrec(precision).Quo(xPowerI, new(big.Float).SetInt(factorialI))
		if i/2%2 != 0 { // Subtract term if i is odd
			result.Sub(result, term)
		} else { // Add term if i is even
			result.Add(result, term)
		}

		// Prepare next term
		xPowerI.Mul(xPowerI, x).Mul(xPowerI, x)                                      // x^(i+2)
		factorialI.Mul(factorialI, big.NewInt(i+1)).Mul(factorialI, big.NewInt(i+2)) // (i+1)! * (i+2)!

		// Break if the term is sufficiently small
		if term.Abs(term).Cmp(new(big.Float).SetPrec(precision).SetFloat64(1e-10)) <= 0 {
			break
		}
	}

	return result, nil
}

func Cos(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cos requires exactly one argument")
	}
	x := args[0]
	result := new(big.Float).SetPrec(precision).SetFloat64(1) // Initialize result with the first term of the series
	xSquared := new(big.Float).SetPrec(precision).Mul(x, x)   // x^2 to be used in the loop
	xTerm := new(big.Float).SetPrec(precision).SetFloat64(1)  // x^0 = 1 for the first term
	sign := int64(-1)                                         // Sign starts as negative for the second term

	for i := int64(2); ; i += 2 {
		xTerm.Mul(xTerm, xSquared) // xTerm *= x^2 to get x^i
		factI := new(big.Float).SetPrec(precision).SetInt(factorial(big.NewInt(i)))
		term := new(big.Float).SetPrec(precision).Quo(xTerm, factI) // (x^i) / i!

		if sign == -1 {
			result.Sub(result, term)
		} else {
			result.Add(result, term)
		}
		sign *= -1 // Alternate sign

		// Break if the term is sufficiently small to ensure convergence
		if term.Abs(term).Cmp(new(big.Float).SetPrec(precision).SetFloat64(1e-50)) <= 0 {
			break
		}
	}

	return result, nil
}

func Tan(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("tan requires exactly one argument")
	}
	sine, err := Sin(precision, args[0])
	if err != nil {
		return nil, err
	}
	cosine, err := Cos(precision, args[0])
	if err != nil {
		return nil, err
	}
	result := new(big.Float).SetPrec(precision).Quo(sine, cosine) // tan(x) = sin(x) / cos(x)
	return result, nil
}

// Asin calculates the arc-sine of a number
func Asin(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("asin requires exactly one argument")
	}

	// This uses the math.Asin function for demonstration, which works on float64.
	// For high precision, a numerical method tailored to *big.Float should be used.
	argFloat64, _ := args[0].Float64()
	asinResult := math.Asin(argFloat64)

	return new(big.Float).SetPrec(precision).SetFloat64(asinResult), nil
}

// Acos calculates the arc cosine of a number.
func Acos(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("acos requires exactly one argument")
	}

	argFloat64, _ := args[0].Float64()
	acosResult := math.Acos(argFloat64) // Using math.Acos for demonstration

	return new(big.Float).SetPrec(precision).SetFloat64(acosResult), nil
}
