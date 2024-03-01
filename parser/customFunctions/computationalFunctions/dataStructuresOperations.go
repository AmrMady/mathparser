package computationalFunctions

import (
	"fmt"
	"math/big"
)

// Assuming Stack is a type you've defined that implements a stack with *big.Float
type Stack struct {
	items []*big.Float
}

// Push adds an item to the stack.
func (s *Stack) Push(item *big.Float) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack.
func (s *Stack) Pop() (*big.Float, error) {
	if len(s.items) == 0 {
		return nil, fmt.Errorf("pop from empty stack")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, nil
}
