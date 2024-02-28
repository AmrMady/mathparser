package parser

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
)

const (
	precision uint = 512
)

func ParseSimple(expression string) (float64, error) {
	tokens := tokenize(expression)
	outputQueue, err := shuntingYard(tokens)
	if err != nil {
		return 0, err
	}

	result, err := evaluateRPN(outputQueue)
	if err != nil {
		return 0, err
	}
	result64, _ := result.Float64()

	return result64, nil
}

func ParseWithVariables(expression string, variables map[string]float64) (float64, error) {
	substitutedExpression := substituteVariables(expression, variables)

	return ParseSimple(substitutedExpression)
}

func substituteVariables(expression string, variables map[string]float64) string {
	for varName, varValue := range variables {
		expression = strings.Replace(expression, varName, strconv.FormatFloat(varValue, 'f', -1, 64), -1)
	}
	return expression
}

func tokenize(expression string) []string {
	var tokens []string
	var currentToken strings.Builder
	previousChar := ' ' // Initialize with a space to handle unary minus at the beginning

	for i, char := range expression {
		if unicode.IsDigit(char) || (char == '.' && i > 0 && unicode.IsDigit(rune(expression[i-1]))) {
			currentToken.WriteRune(char)
		} else if char == '-' && (i == 0 || previousChar == '(' || isOperator(previousChar, i-1, expression)) {
			// Check if this '-' is unary
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			currentToken.WriteRune(char) // Start a new token with '-'
		} else {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			if char != ' ' {
				tokens = append(tokens, string(char))
			}
		}
		if char != ' ' {
			previousChar = char
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func isOperator(char rune, charIndex int, fullString string) bool {
	operators := "+-*/^"
	return strings.ContainsRune(operators, char) && char != '-' || (char == '-' && charIndex > 0 && !unicode.IsDigit(rune(fullString[charIndex-1])) && fullString[charIndex-1] != ')')
}

func isDigit(char rune) bool {
	return (char >= '0' && char <= '9') || char == '.'
}

func shuntingYard(tokens []string) ([]string, error) {
	var outputQueue []string
	var operatorStack []string

	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2, "^": 3}

	for _, token := range tokens {
		switch {
		case isDigit(rune(token[0])) || (token[0] == '-' && len(token) > 1): // Number or unary minus with a number
			outputQueue = append(outputQueue, token)
		case token == "(":
			operatorStack = append(operatorStack, token)
		case token == ")":
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != "(" {
				outputQueue = append(outputQueue, operatorStack[len(operatorStack)-1])
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			operatorStack = operatorStack[:len(operatorStack)-1] // Pop the "("
		default: // Operator
			for len(operatorStack) > 0 && precedence[operatorStack[len(operatorStack)-1]] >= precedence[token] {
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

func evaluateRPN(tokens []string) (*big.Float, error) {
	var stack []*big.Float

	for _, token := range tokens {
		if isDigit(rune(token[0])) || (token[0] == '-' && len(token) > 1) { // Check for number or unary minus
			num, _, err := new(big.Float).SetPrec(precision).Parse(token, 10)
			if err != nil {
				return nil, err
			}
			stack = append(stack, num)
		} else { // Operator
			if len(stack) < 2 && token != "-" { // Unary minus is an exception, can operate with one operand
				return nil, fmt.Errorf("insufficient values in the stack for operation %s", token)
			}

			// Pop operands from the stack
			var b, a *big.Float                  // Initialize operands
			if token == "-" && len(stack) == 1 { // Unary minus case
				a = stack[len(stack)-1]
				stack = stack[:len(stack)-1]    // Pop one operand
				result := new(big.Float).Neg(a) // Unary minus operation
				stack = append(stack, result)
			} else { // Binary operations
				b = stack[len(stack)-1]
				a = stack[len(stack)-2]
				stack = stack[:len(stack)-2] // Pop two operands

				var result *big.Float
				switch token {
				case "+":
					result = new(big.Float).SetPrec(precision).Add(a, b)
				case "-":
					result = new(big.Float).SetPrec(precision).Sub(a, b)
				case "*":
					result = new(big.Float).SetPrec(precision).Mul(a, b)
				case "/":
					if b.Cmp(new(big.Float).SetPrec(precision).SetFloat64(0)) == 0 {
						return nil, fmt.Errorf("division by zero")
					}
					result = new(big.Float).Quo(a, b)
				case "^":
					// Assume binaryExponentiation function is correctly implemented to handle big.Float
					result = binaryExponentiation(a, b)
				}
				stack = append(stack, result)
			}
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("evaluation error: stack has unexpected size")
	}
	return stack[0], nil
}

func binaryExponentiation(base, exponent *big.Float) *big.Float {
	expInt, _ := exponent.Int64()
	result := new(big.Float).SetPrec(precision).SetFloat64(1)
	for expInt > 0 {
		if expInt&1 == 1 {
			result.SetPrec(precision).Mul(result, base)
		}
		base.SetPrec(precision).Mul(base, base)
		expInt >>= 1
	}
	return result
}
