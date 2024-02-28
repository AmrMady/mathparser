package parser

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
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
	for varName := range variables {
		placeholder := fmt.Sprintf("${%s}", varName)
		expression = strings.Replace(expression, varName, placeholder, -1)
	}
	for varName, varValue := range variables {
		placeholder := fmt.Sprintf("${%s}", varName)
		formattedValue := strconv.FormatFloat(varValue, 'f', -1, 64)
		expression = strings.Replace(expression, placeholder, formattedValue, -1)
	}
	return expression
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

func evaluateRPN(tokens []string) (*big.Float, error) {
	var stack []*big.Float

	for _, token := range tokens {
		if isDigit(rune(token[0])) {
			fmt.Println("token: ", token)
			num, _, err := new(big.Float).SetPrec(precision).Parse(token, 10)
			if err != nil {
				return nil, err
			}
			fmt.Println("num: ", num)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in the stack for operation %s", token)
			}
			a, b := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			var result *big.Float
			switch token {
			case "+":
				result = new(big.Float).SetPrec(precision).Add(a, b)
			case "-":
				result = new(big.Float).SetPrec(precision).Sub(a, b)
			case "*":
				result = new(big.Float).SetPrec(precision).Mul(a, b)
			case "/":
				if b.Cmp(new(big.Float).SetFloat64(0)) == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				result = new(big.Float).SetPrec(precision).Quo(a, b)
			case "^":
				result = binaryExponentiation(a, b)
			}
			stack = append(stack, result)
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
