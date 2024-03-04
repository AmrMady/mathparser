package basicMathematicalFunctions

import (
	"fmt"
	"math/big"
)

// Pow raises a number to the power of another
func Pow(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("pow requires exactly two arguments")
	}
	base := args[0].SetPrec(precision)
	exponent := args[1].SetPrec(precision)
	result, err := BinaryExponentiation(precision, base, exponent) // Ensure binaryExponentiation can handle *big.Float
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Sqrt calculates the square root of a number
func Sqrt(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("sqrt requires exactly one argument")
	}
	x := args[0].SetPrec(precision)
	result := new(big.Float).SetPrec(precision)
	result.Sqrt(x)
	return result, nil
}

// Cbrt calculates the cube root of a number
func Cbrt(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("cbrt requires exactly one argument")
	}
	x := args[0]

	// Cube root calculation using Newton's method for demonstration
	result := new(big.Float).SetPrec(precision).Set(x)
	three := new(big.Float).SetPrec(precision).SetFloat64(3)
	two := new(big.Float).SetPrec(precision).SetFloat64(2)
	for i := 0; i < 100; i++ {
		// Newton's iteration: result = (2*result + x/(result*result)) / 3
		temp := new(big.Float).Quo(x, new(big.Float).Mul(result, result))
		temp.Add(temp, new(big.Float).Mul(two, result))
		result.Quo(temp, three)
	}
	return result, nil
}
