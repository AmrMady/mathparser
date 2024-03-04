package basicMathematicalFunctions

import (
	"math"
	"math/big"
	"testing"
)

func TestTan(t *testing.T) {
	precision := uint(512)
	pi := big.NewFloat(math.Pi)
	four := big.NewFloat(4)
	piOverFour := new(big.Float).Quo(pi, four)

	expected := big.NewFloat(1) // Expected result for tan(pi/4)
	result, err := Tan(precision, piOverFour)
	if err != nil {
		t.Errorf("Tan function returned an error: %v", err)
	}

	// Compare the result with the expected value
	if result.Cmp(expected) != 0 {
		t.Errorf("Expected Tan(pi/4) = 1, got %v", result)
	}
}
