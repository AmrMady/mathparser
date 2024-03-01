package computationalFunctions

import (
	"fmt"
	"math/big"
)

// Sort sorts a slice of *big.Float numbers
func Sort(precision uint, args ...*big.Float) ([]*big.Float, error) {
	if len(args) < 2 {
		return args, nil // No sorting needed for less than 2 elements
	}

	// Bubble Sort for demonstration
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(args)-i-1; j++ {
			if args[j].Cmp(args[j+1]) > 0 {
				args[j], args[j+1] = args[j+1], args[j]
			}
		}
	}
	return args, nil
}

// BinarySearch finds the position of a target value within a sorted array.
func BinarySearch(precision uint, args ...*big.Float) (*big.Float, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("binarySearch requires at least two arguments: the array (sorted) and the target value")
	}
	target := args[len(args)-1]       // Last argument is the target value
	sortedArray := args[:len(args)-1] // All but the last argument

	low := 0
	high := len(sortedArray) - 1

	for low <= high {
		mid := (low + high) / 2
		compareResult := sortedArray[mid].Cmp(target)

		if compareResult == 0 {
			return new(big.Float).SetPrec(precision).SetInt64(int64(mid)), nil // Target found
		} else if compareResult < 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return new(big.Float).SetPrec(precision).SetInt64(-1), nil // Target not found
}
