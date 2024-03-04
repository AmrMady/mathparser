package advancedMathematicalFunctions

import (
	"fmt"
	"math/big"
)

// InsertionSort for sorting slice of *big.Float
func insertionSort(values []*big.Float) {
	for i := 1; i < len(values); i++ {
		key := values[i]
		j := i - 1

		// Move elements of values[0..i-1], that are greater than key, to one position ahead of their current position
		for j >= 0 && values[j].Cmp(key) > 0 {
			values[j+1] = values[j]
			j = j - 1
		}
		values[j+1] = key
	}
}

func Mean(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("mean requires at least one argument")
	}

	sum := new(big.Float).SetPrec(precision).SetFloat64(0)
	for _, arg := range args {
		sum.Add(sum, arg)
	}

	count := new(big.Float).SetPrec(precision).SetInt64(int64(len(args)))
	mean := new(big.Float).Quo(sum, count)
	return mean, nil
}

// StdDev calculates the standard deviation of a given list of numbers
func StdDev(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("StdDev requires at least two arguments")
	}
	n := len(args)
	mean, err := Mean(precision, args...)
	if err != nil {
		return nil, err
	}

	// Calculate the variance
	varSum := new(big.Float).SetPrec(precision).SetFloat64(0)
	for _, val := range args {
		diff := new(big.Float).Sub(val, mean)
		sqDiff := new(big.Float).Mul(diff, diff)
		varSum.Add(varSum, sqDiff)
	}
	nBig := new(big.Float).SetPrec(precision).SetFloat64(float64(n))
	variance := new(big.Float).Quo(varSum, nBig)

	// Standard deviation is the square root of variance
	stdDev := new(big.Float).SetPrec(precision)
	stdDev.Sqrt(variance) // sqrt(variance)
	return stdDev, nil
}

// Median calculates the middle value of a given list of numbers
func Median(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("median requires at least one argument")
	}

	// Create a copy of args to sort, to not modify the original slice
	values := make([]*big.Float, len(args))
	copy(values, args)

	// Sort the copy of args
	insertionSort(values)

	// Calculate median
	n := len(values)
	half := n / 2
	median := new(big.Float).SetPrec(precision)

	if n%2 == 0 {
		// even number of elements, average the two middle
		median.Add(values[half-1], values[half]).Quo(median, big.NewFloat(2))
	} else {
		// odd number of elements, take the middle
		median = values[half]
	}

	return median, nil
}
