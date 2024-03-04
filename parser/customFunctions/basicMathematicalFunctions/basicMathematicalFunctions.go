package basicMathematicalFunctions

import (
	"math/big"
)

func factorial(n *big.Int) *big.Int {
	one := big.NewInt(1)
	if n.Cmp(one) == 0 || n.Cmp(big.NewInt(0)) == 0 {
		return one
	}
	return new(big.Int).Mul(n, factorial(new(big.Int).Sub(n, one)))
}

func BinaryExponentiation(precision uint, base, exponent *big.Float) (*big.Float, error) {
	// Check if exponent is an integer
	exponentInt, accuracy := exponent.Int64()
	if accuracy == big.Exact { // The exponent is an integer
		result := new(big.Float).SetPrec(precision).SetFloat64(1)
		for exponentInt > 0 {
			if exponentInt&1 == 1 {
				result.SetPrec(precision).Mul(result, base)
			}
			base.Mul(base, base)
			exponentInt >>= 1
		}
		return result, nil
	} else { // The exponent is not an integer, use a^b = e^(b * ln(a))
		lnBase, err := Log(precision, base) // Calculate ln(base)
		if err != nil {
			return nil, err
		}
		exponentTimesLnBase := new(big.Float).SetPrec(precision).Mul(exponent, lnBase) // Calculate exponent * ln(base)
		result, err := Exp(precision, exponentTimesLnBase)                             // Calculate e^(exponent * ln(base))
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}
