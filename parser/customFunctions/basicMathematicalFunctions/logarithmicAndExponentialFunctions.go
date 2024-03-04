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
	if x.Cmp(big.NewFloat(1)) == 0 {
		return big.NewFloat(0), nil
	}

	// Ensure the precision is set for all calculations
	x = new(big.Float).SetPrec(precision).Set(x)
	one := new(big.Float).SetPrec(precision).SetFloat64(1)

	// Improved accuracy for (x-1)/(x+1)
	xMinusOne := new(big.Float).Sub(x, one)
	xPlusOne := new(big.Float).Add(x, one)
	fraction := new(big.Float).Quo(xMinusOne, xPlusOne)

	result := new(big.Float).SetPrec(precision)
	term := new(big.Float).SetPrec(precision)
	fractionSquared := new(big.Float).Mul(fraction, fraction)

	// Using a more efficient series calculation
	for n := int64(1); ; n += 2 {
		var err error
		exponent := new(big.Float).SetInt64(n)
		// Calculate (fraction^2)^(n) directly without using BinaryExponentiation for each term
		term, err = Exp(precision, fractionSquared, exponent)
		if err != nil {
			return big.NewFloat(0), nil
		}
		term.Quo(term, new(big.Float).SetPrec(precision).SetFloat64(float64(n))) // term/n
		if n%4 == 1 {
			result.Add(result, term)
		} else {
			result.Sub(result, term)
		}

		// Break if the term is sufficiently small to ensure convergence
		if term.Abs(term).Cmp(new(big.Float).SetPrec(precision).SetFloat64(1e-10)) <= 0 {
			break
		}
	}
	result.Mul(result, big.NewFloat(2))                   // Multiply the sum by 2
	result.Add(result, new(big.Float).Quo(fraction, one)) // Add the first term of the series

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
		term.SetPrec(precision).Mul(term, x)                                                      // x^i
		term.SetPrec(precision).Quo(term, new(big.Float).SetInt(factorial(big.NewInt(int64(i))))) // x^i / i!
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
