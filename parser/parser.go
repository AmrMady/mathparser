package parser

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	precision uint = 512
)

func Parse(expression string) (float64, error) {
	tokens := tokenize(expression)
	outputQueue, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evaluateRPN(outputQueue)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func tokenize(expression string) []string {
	var tokens []string
	currentToken := ""

	for _, char := range expression {
		if isOperator(char) {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		} else if isDigit(char) {
			currentToken += string(char)
		} else if char == '(' || char == ')' {
			if currentToken != "" {
				tokens = append(tokens, currentToken)
				currentToken = ""
			}
			tokens = append(tokens, string(char))
		}
	}

	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

func isOperator(char rune) bool {
	operators := "+-*/^"
	return strings.ContainsRune(operators, char)
}

func isDigit(char rune) bool {
	return (char >= '0' && char <= '9') || char == '.'
}

func shuntingYard(tokens []string) ([]string, error) {
	var outputQueue []string
	var operatorStack []string

	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3}

	for _, token := range tokens {
		if isDigit(rune(token[0])) {
			outputQueue = append(outputQueue, token)
		} else if token == "(" {
			operatorStack = append(operatorStack, token)
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			operatorStack = operatorStack[:len(operatorStack)-1]
		} else if isOperator(rune(token[0])) {
			for len(operatorStack) > 0 && (precedence[operatorStack[len(operatorStack)-1]] > precedence[token]) {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, token)
		}
	}

	for len(operatorStack) > 0 {
		if operatorStack[len(operatorStack)-1] == "(" || operatorStack[len(operatorStack)-1] == ")" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return outputQueue, nil
}
func evaluateRPN(tokens []string) (float64, error) {
	var stack []*big.Float

	for _, token := range tokens {
		if isDigit(rune(token[0])) {
			f := new(big.Float).SetPrec(precision)
			bigFloat, _, err := f.Parse(token, 10)
			if err != nil {
				fmt.Println("Error:", err)
				return 0, err
			}
			stack = append(stack, bigFloat)
		} else if isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, fmt.Errorf("insufficient operands for operator %s", token)
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result big.Float
			switch token {
			case "+":
				result.Add(a, b)
			case "-":
				result.Sub(a, b)
			case "*":
				result.Mul(a, b)
			case "/":
				if b.Sign() == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				result.Quo(a, b)
			case "^":
				result = *new(big.Float).SetPrec(precision).SetFloat64(1)
				exponent, _ := b.Int64()
				for i := int64(0); i < exponent; i++ {
					result.Mul(&result, a)
				}

			}
			stack = append(stack, &result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}

	result64, _ := stack[0].Float64()

	return result64, nil
}
