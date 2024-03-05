package main

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser"
	"github.com/AmrMady/mathparser/parser/customFunctions/basicMathematicalFunctions"
	"math/big"
)

func main() {

	result, err := parser.ParseAndReturnString("sin(pi/6) + log(e^2) + (sqrt(16) * 3)")
	if err != nil {
		fmt.Println("ParseAndReturnString::Error:", err)
		return
	}

	fmt.Println("result as String:", result)

}

//sin(pi/6) + log(e^2) + sqrt(16)*3
//exp(1) * (cos(0) + tan(pi/4))
//"exp(1) * (cos(0) + tan(pi/4))"

//expression := "x * y + (z - w) / a"
//
//// Define variables and their values
//variables := map[string]float64{
//	"x": 5.5,
//	"y": 4,
//	"z": 20,
//	"w": 15,
//	"a": 2.5,
//}
//
//// Parse and evaluate the expression with variables
//result, err := parser.ParseWithVariables(expression, variables)
//if err != nil {
//	fmt.Println("Error:", err)
//	return
//}
//fmt.Printf("Result of '%s': %f\n", expression, result)

func CalculatePiBBP(iterations int) (*big.Float, error) {
	precision := uint(512)
	pi := new(big.Float).SetPrec(precision).SetFloat64(0)
	one := new(big.Float).SetPrec(precision).SetFloat64(1)
	four := new(big.Float).SetPrec(precision).SetFloat64(4)
	two := new(big.Float).SetPrec(precision).SetFloat64(2)
	sixteen := new(big.Float).SetPrec(precision).SetFloat64(16)

	for k := 0; k < iterations; k++ {
		kBig := new(big.Float).SetPrec(precision).SetInt64(int64(k))
		eightK := new(big.Float).SetPrec(precision).Mul(new(big.Float).SetPrec(precision).SetInt64(8), kBig)

		fraction1 := new(big.Float).Quo(four, new(big.Float).Add(eightK, one))
		fraction2 := new(big.Float).Quo(two, new(big.Float).Add(eightK, new(big.Float).SetInt64(4)))
		fraction3 := new(big.Float).Quo(one, new(big.Float).Add(eightK, new(big.Float).SetInt64(5)))
		fraction4 := new(big.Float).Quo(one, new(big.Float).Add(eightK, new(big.Float).SetInt64(6)))

		// Calculate 16^-k
		divisor, err := basicMathematicalFunctions.BinaryExponentiation(512, sixteen, kBig)
		if err != nil {
			fmt.Println("CalculatePiBBP::BinaryExponentiation::Error:", err)
			return nil, err
		}
		sixteenToNegK := new(big.Float).Quo(one, divisor)

		// Sum the fractions and multiply by 16^-k
		sumFractions := new(big.Float).Sub(fraction1, fraction2)
		sumFractions.Sub(sumFractions, fraction3)
		sumFractions.Sub(sumFractions, fraction4)
		sumFractions.Mul(sumFractions, sixteenToNegK)

		pi.Add(pi, sumFractions)
	}

	return pi, nil
}

func main2() {
	res, err := CalculatePiBBP(1000)
	if err != nil {
		fmt.Println("CalculatePiBBP::Error:", err)
		return
	}
	fmt.Println("Pi as String:", res.Text('f', -1))
}
