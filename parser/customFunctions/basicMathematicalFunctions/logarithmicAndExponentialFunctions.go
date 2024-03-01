package basicMathematicalFunctions

import (
	"fmt"
	"math/big"
)

// Log calculates natural logarithm using Taylor series approximation
func Log(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log requires exactly one argument")
	}
	x := args[0]
	if x.SetPrec(precision).Cmp(big.NewFloat(1)) == 0 {
		return big.NewFloat(0), nil
	}

	result := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision)
	one := big.NewFloat(1).SetPrec(precision)

	// (x-1)/(x+1)
	xMinusOne := new(big.Float).SetPrec(precision).Sub(x, one)
	xPlusOne := new(big.Float).SetPrec(precision).Add(x, one)
	fraction := new(big.Float).SetPrec(precision).Quo(xMinusOne, xPlusOne)

	fractionSquared := new(big.Float).SetPrec(precision).Mul(fraction, fraction)

	for i := int64(1); i < 100; i += 2 {
		// Calculate each term
		var err error
		exponent := new(big.Float).SetInt64(i).SetPrec(precision)
		term, err = BinaryExponentiation(precision, fractionSquared, exponent) // (fraction^2)^n
		if err != nil {
			return nil, err
		}
		term.SetPrec(precision).Quo(term, big.NewFloat(float64(i)).SetPrec(precision)) // term / n
		result.Add(result, term)
	}
	result.SetPrec(precision).Mul(result, big.NewFloat(2)) // Multiply the sum by 2
	result.SetPrec(precision).Add(result, fraction)        // Add the first term of the series

	return result, nil
}

func Exp(precision uint, args ...*big.Float) (*big.Float, error) { // exponential function || e^x
	if len(args) != 1 {
		return nil, fmt.Errorf("exp requires exactly one argument")
	}
	x := args[0]
	result := new(big.Float).SetPrec(precision).SetFloat64(1) // Start with 1
	term := new(big.Float).SetPrec(precision).SetFloat64(1)

	for i := 1; i < 20; i++ {
		term.SetPrec(precision).Mul(term, x)                                                             // x^i
		term.SetPrec(precision).Quo(term, new(big.Float).SetPrec(precision).SetInt(factorial(int64(i)))) // x^i / i!
		result.Add(result, term)
	}

	return result, nil
}

// Log10 calculates the base-10 logarithm of a number
func Log10(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("log10 requires exactly one argument")
	}
	err := error(nil)
	x := args[0].SetPrec(precision)
	log := new(big.Float).SetPrec(precision)
	// Convert ln(x) to log10(x) = ln(x) / ln(10)
	ln10 := new(big.Float).SetPrec(precision).SetFloat64(2.302585092994046) // Precomputed ln(10)
	log, err = Log(precision, x)                                            // Natural log of x
	if err != nil {
		return nil, err
	}
	log.Quo(log, ln10) // Divide by ln(10) to convert to base-10
	return log, nil
}
