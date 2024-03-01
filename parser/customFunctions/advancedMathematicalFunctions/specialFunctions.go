package advancedMathematicalFunctions

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser/customFunctions/basicMathematicalFunctions"
	"math/big"
)

// Gamma approximates the Gamma function using Stirling's approximation.
func Gamma(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("gamma requires exactly one argument")
	}
	x := args[0]
	if x.Cmp(big.NewFloat(1)) == 0 || x.Cmp(big.NewFloat(2)) == 0 {
		return big.NewFloat(1).SetPrec(precision), nil // Gamma(n) = (n-1)!
	}

	// Stirling's approximation: Gamma(z) ~ sqrt(2*pi/z) * (z/e)^z
	twoPi := new(big.Float).Mul(big.NewFloat(2), big.NewFloat(3.141592653589793238462643383279)).SetPrec(precision)
	sqrtTwoPiOverZ := new(big.Float).Quo(twoPi, x)
	sqrtTwoPiOverZ.Sqrt(sqrtTwoPiOverZ)

	zOverE := new(big.Float).Quo(x, big.NewFloat(2.7182818284590452353602874713527)).SetPrec(precision)
	zOverEPowZ, err := basicMathematicalFunctions.BinaryExponentiation(precision, zOverE, x)
	if err != nil {
		return nil, err
	}

	gamma := new(big.Float).Mul(sqrtTwoPiOverZ, zOverEPowZ).SetPrec(precision)
	return gamma, nil
}
