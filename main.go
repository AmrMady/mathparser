package main

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser"
	"github.com/AmrMady/mathparser/parser/customFunctions/basicMathematicalFunctions"
	"math/big"
)

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

func main() {
	piString, err := parser.ParseAndReturnString("1 / (12 * (sum of ((-1)^k * (6k)! * (545140134k + 13591409) / ((3k)! * (k!)^3 * 640320^(3k + 1.5)) from k=0 to infinity))")
	if err != nil {
		fmt.Println("ParseAndReturnString::Error:", err)
		return
	}

	fmt.Println("Pi as String:", piString)

	//res, err := CalculatePiBBP(1000)
	//if err != nil {
	//	fmt.Println("CalculatePiBBP::Error:", err)
	//	return
	//}
	//fmt.Println("Pi as String:", res.Text('f', -1))
}
