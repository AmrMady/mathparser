package customFunctions

import (
	"errors"
	"github.com/AmrMady/mathparser/parser/customFunctions/basicMathematicalFunctions"
	"math/big"
)

func BbpTerm(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) != 1 {
		return nil, errors.New("bbpTerm requires exactly one argument")
	}

	k := args[0]
	four := big.NewFloat(4).SetPrec(precision)
	two := big.NewFloat(2).SetPrec(precision)
	one := big.NewFloat(1).SetPrec(precision)
	sixteen := big.NewFloat(16).SetPrec(precision)

	// Calculating 16^-k part of the term
	divisor, err := basicMathematicalFunctions.BinaryExponentiation(precision, sixteen, k)
	if err != nil {
		return nil, err
	}
	sixteenNegK := new(big.Float).Quo(one, divisor)

	// Calculating the series part of the term
	seriesPart := new(big.Float).Quo(four, new(big.Float).Add(new(big.Float).Mul(big.NewFloat(8), k), one))
	seriesPart.Sub(seriesPart, new(big.Float).Quo(two, new(big.Float).Add(new(big.Float).Mul(big.NewFloat(8), k), big.NewFloat(4))))
	seriesPart.Sub(seriesPart, new(big.Float).Quo(one, new(big.Float).Add(new(big.Float).Mul(big.NewFloat(8), k), big.NewFloat(5))))
	seriesPart.Sub(seriesPart, new(big.Float).Quo(one, new(big.Float).Add(new(big.Float).Mul(big.NewFloat(8), k), big.NewFloat(6))))

	// Combine parts to get the term result
	termResult := new(big.Float).Mul(seriesPart, sixteenNegK).SetPrec(precision)

	return termResult, nil
}
