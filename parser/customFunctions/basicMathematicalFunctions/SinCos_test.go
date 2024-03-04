package basicMathematicalFunctions

import (
	"math"
	"math/big"
	"testing"
)

func TestSinCos(t *testing.T) {
	precision := uint(512)
	pi := big.NewFloat(math.Pi)
	four := big.NewFloat(4)
	piOverFour := new(big.Float).Quo(pi, four)

	expectedSin := big.NewFloat(math.Sqrt(2) / 2) // sin(pi/4) = sqrt(2)/2
	resultSin, _ := Sin(precision, piOverFour)

	if resultSin.Cmp(expectedSin) != 0 {
		t.Errorf("Expected Sin(pi/4) = %v, got %v", expectedSin, resultSin)
	}

	expectedCos := big.NewFloat(math.Sqrt(2) / 2) // cos(pi/4) = sqrt(2)/2
	resultCos, _ := Cos(precision, piOverFour)

	if resultCos.Cmp(expectedCos) != 0 {
		t.Errorf("Expected Cos(pi/4) = %v, got %v", expectedCos, resultCos)
	}
}
