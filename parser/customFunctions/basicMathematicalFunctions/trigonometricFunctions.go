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
	term := new(big.Float).SetPrec(precision)

	for i, sign := 0, int64(1); i < 20; i, sign = i+1, -sign {
		n := 2*i + 1
		nBig := new(big.Float).SetPrec(precision).SetInt64(int64(n))
		factorialN := new(big.Float).SetInt(factorial(int64(n)))
		var err error
		term, err = BinaryExponentiation(precision, x, nBig) // x^n
		if err != nil {
			return nil, err
		}
		term.Quo(term, factorialN)                  // x^n / n!
		term.Mul(term, big.NewFloat(float64(sign))) // Apply sign

		result.Add(result, term)

		// Break if the added term is less than a predefined small threshold
		if term.Abs(term).Cmp(big.NewFloat(1e-50)) == -1 {
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
	result := new(big.Float).SetPrec(precision).SetFloat64(1) // Start with 1
	term := new(big.Float).SetPrec(precision)

	for i, sign := 1, int64(-1); i < 20; i += 2 {
		n := int64(i + 1)
		nBig := new(big.Float).SetPrec(precision).SetInt64(n)
		factorialN := new(big.Float).SetInt(factorial(n))
		var err error
		term, err = BinaryExponentiation(precision, x, nBig) // x^n
		if err != nil {
			return nil, err
		}
		term.Quo(term, factorialN)                  // x^n / n!
		term.Mul(term, big.NewFloat(float64(sign))) // Apply sign

		result.Add(result, term)

		sign = -sign // Flip sign for next term
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
	result := new(big.Float).Quo(sine, cosine) // tan(x) = sin(x) / cos(x)
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
