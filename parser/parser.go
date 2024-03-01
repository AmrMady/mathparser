package parser

import (
	"fmt"
	"github.com/AmrMady/mathparser/parser/customFunctions"
	"github.com/AmrMady/mathparser/parser/customFunctions/advancedMathematicalFunctions"
	"github.com/AmrMady/mathparser/parser/customFunctions/basicMathematicalFunctions"
	"github.com/AmrMady/mathparser/parser/customFunctions/computationalFunctions"
	"math/big"
	"strconv"
	"strings"
)

const (
	precision uint = 512
)

const (
	TOKEN_SUM_OF = iota + 2001
	TOKEN_FROM
	TOKEN_TO
	TOKEN_INFINITY
	TOKEN_IDENTIFIER
	TOKEN_LPAREN
	TOKEN_RPAREN
	TOKEN_COMMA
	TOKEN_NUMBER
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_ASTERISK
	TOKEN_SLASH
	TOKEN_CARET
	TOKEN_EQUALS
	TOKEN_FACTORIAL
	TOKEN_START
	TOKEN_ERROR
	TOKEN_EOF
	// You can add more token types here as needed.
)

var (
	customFuncsMap   = make(map[string]func(precision uint, args ...*big.Float) (*big.Float, error))
	customActionsMap = make(map[string]interface{}) // For operations/actions with different signatures
	variableContext  = map[string]*big.Float{}
)

// RegisterCustomFunction for mathematical operations with a specific signature
func RegisterCustomFunction(name string, function func(precision uint, args ...*big.Float) (*big.Float, error)) {
	customFuncsMap[name] = function
}

// RegisterCustomAction for operations/actions with varied signatures or operational needs
func RegisterCustomAction(name string, action interface{}) {
	customActionsMap[name] = action
}

func init() {
	// Register the bbpTerm custom function
	RegisterCustomFunction("bbpTerm", customFunctions.BbpTerm)
	RegisterCustomFunction("sin", basicMathematicalFunctions.Sin)
	RegisterCustomFunction("cos", basicMathematicalFunctions.Cos)
	RegisterCustomFunction("tan", basicMathematicalFunctions.Tan)
	RegisterCustomFunction("exp", basicMathematicalFunctions.Exp)
	RegisterCustomFunction("log10", basicMathematicalFunctions.Log10)
	RegisterCustomFunction("sqrt", basicMathematicalFunctions.Sqrt)
	RegisterCustomFunction("pow", basicMathematicalFunctions.Pow)

	RegisterCustomFunction("log", basicMathematicalFunctions.Log)

	RegisterCustomFunction("gamma", advancedMathematicalFunctions.Gamma)
	RegisterCustomFunction("median", advancedMathematicalFunctions.Median)
	RegisterCustomFunction("stddev", advancedMathematicalFunctions.StdDev)

	RegisterCustomFunction("stddev", advancedMathematicalFunctions.StdDev)
	RegisterCustomFunction("isPrime", advancedMathematicalFunctions.IsPrime)
	RegisterCustomFunction("gcd", advancedMathematicalFunctions.GCD)
	RegisterCustomFunction("mean", advancedMathematicalFunctions.Mean)

	RegisterCustomFunction("asin", basicMathematicalFunctions.Asin)
	RegisterCustomFunction("cbrt", basicMathematicalFunctions.Cbrt)
	RegisterCustomAction("sort", computationalFunctions.Sort)
}

// Parser constructs an AST from tokens.
type Parser struct {
	lexer           *Lexer
	currentToken    InputToken
	prevToken       *InputToken
	variableContext map[string]*big.Float
}

// NewParser creates a new parser instance.
func NewParser(lexer *Lexer) *Parser {
	return &Parser{lexer: lexer}
}

func (p *Parser) setVariable(name string, value *big.Float) {
	if p.variableContext == nil {
		p.variableContext = make(map[string]*big.Float)
	}
	p.variableContext[name] = value
}

func (p *Parser) advance() {
	if p.currentToken.Type != TOKEN_EOF {
		p.prevToken = &p.currentToken
	}
	p.currentToken = p.lexer.NextToken()
}

func (p *Parser) expect(tokenType int) bool {
	if p.currentToken.Type == tokenType {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) isUnaryMinusContext() bool {
	// True if it's the start of the expression.
	if p.prevToken == nil {
		return true
	}
	// True if the previous token was an operator or an opening parenthesis.
	switch p.prevToken.Type {
	case TOKEN_LPAREN, TOKEN_PLUS, TOKEN_MINUS, TOKEN_ASTERISK, TOKEN_SLASH, TOKEN_CARET:
		return true
	}
	return false
}

func (p *Parser) parseExpression() (Node, error) {
	// First, check if the expression starts with a summation
	if p.currentToken.Type == TOKEN_SUM_OF {
		return p.parseSummationExpression()
	}

	// Parse the first term
	left, err := p.parseTerm()
	if err != nil {
		fmt.Println("parseExpression::parseTerm::left error: ", err)
		return nil, err
	}

	// Handle binary operations
	for p.isOperatorToken(p.currentToken.Type) {
		op := p.currentToken // Store the current operator token
		p.advance()          // Move past the operator

		right, err := p.parseTerm() // Parse the next term
		if err != nil {
			fmt.Println("parseExpression::parseTerm::right error: ", err)
			return nil, err
		}

		// Construct a binary operation node from the left-hand side, operator, and right-hand side
		left = &BinaryOpNode{
			Left:  left,
			Op:    op.Value,
			Right: right,
		}
	}

	return left, nil
}

func (p *Parser) peekToken() InputToken {
	// Save the current lexer position.
	savedPos := p.lexer.pos

	// Fetch the next token using the lexer.
	nextToken := p.lexer.NextToken()

	// Restore the lexer position.
	p.lexer.pos = savedPos

	return nextToken
}

func (p *Parser) parseSumSeries() Node {
	p.expect(TOKEN_SUM_OF) // Consume 'sum of'

	// Parse the series term, variable, start, and end expressions
	term, _ := p.parseTerm() // Implement parseTerm based on your needs

	p.expect(TOKEN_FROM)
	varName := p.currentToken.Value // Assuming variable name follows 'from'
	p.advance()

	start, _ := p.parseExpression() // Implement parseExpression for general expression parsing
	p.expect(TOKEN_TO)
	end, _ := p.parseExpression()

	return &SumSeriesNode{
		Term:    term,
		VarName: varName,
		Start:   start,
		End:     end,
	}
}

func (p *Parser) parseTerm() (Node, error) {
	// First handle any unary minus cases
	if p.currentToken.Type == TOKEN_MINUS && p.isUnaryMinusContext() {
		p.advance()                      // Consume the unary minus
		operand, err := p.parsePrimary() // Now directly parsing primary to handle factorials
		if err != nil {
			return nil, err
		}
		return &UnaryOpNode{Op: "-", Operand: operand}, nil
	}

	// For other cases, directly parse primary expressions (which includes factorial handling)
	return p.parsePrimary()
}

func (p *Parser) parseFunctionCall(funcName string) (Node, error) {
	p.advance() // Consume the '('

	var args []Node
	if p.currentToken.Type != TOKEN_RPAREN {
		for p.currentToken.Type != TOKEN_RPAREN && p.currentToken.Type != TOKEN_EOF {
			arg, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if p.currentToken.Type == TOKEN_COMMA {
				p.advance() // Skip over the comma to the next argument
			}
		}
	}

	if !p.expect(TOKEN_RPAREN) {
		return nil, fmt.Errorf("expected closing parenthesis for function call, got %v", p.currentToken)
	}

	return &FunctionCallNode{Name: funcName, Arguments: args}, nil
}

// Improved version of parseSummationExpression to handle 'infinity' and variable substitutions.
func (p *Parser) parseSummationExpression() (Node, error) {
	p.advance() // Move past "sum of"

	// Parse the expression part of the summation
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Ensure "from" keyword is present
	if !p.expect(TOKEN_FROM) {
		return nil, fmt.Errorf("expected 'from', found %v", p.currentToken)
	}

	// Parse the variable name
	variable := p.currentToken.Value
	p.advance()

	// Expect and consume the "=" token after the variable name
	if !p.expect(TOKEN_EQUALS) {
		return nil, fmt.Errorf("expected '=', found %v", p.currentToken)
	}

	// Parse the start expression
	start, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Ensure "to" keyword is present
	if !p.expect(TOKEN_TO) {
		return nil, fmt.Errorf("expected 'to', found %v", p.currentToken)
	}

	// Parse the end expression, handling "infinity"
	var end Node
	if p.currentToken.Type == TOKEN_INFINITY {
		end = &InfinityNode{}
		p.advance() // Move past "infinity"
	} else {
		end, err = p.parseExpression()
		if err != nil {
			return nil, err
		}
	}

	// Construct and return a SummationNode
	return &SummationNode{
		Variable:   variable,
		Start:      start,
		End:        end,
		Expression: expression,
	}, nil
}

func (p *Parser) evaluateSummation(sumNode *SummationNode) (*big.Float, error) {
	var sum big.Float
	startVal, err := evaluate(sumNode.Start)
	if err != nil {
		return nil, err
	}
	endVal, err := evaluate(sumNode.End)
	if err != nil {
		return nil, err
	}
	// Assuming startVal and endVal are integers for simplicity
	startValIntPart, _ := startVal.Int64()
	endValIntPart, _ := endVal.Int64()
	for i := startValIntPart; i <= endValIntPart; i++ {
		p.setVariable(sumNode.Variable, new(big.Float).SetInt64(i))
		termVal, err := evaluate(sumNode.Expression)
		if err != nil {
			return nil, err
		}
		sum.Add(&sum, termVal)
	}
	return &sum, nil
}

func (p *Parser) isOperatorToken(tokenType int) bool {
	switch tokenType {
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_ASTERISK, TOKEN_SLASH, TOKEN_CARET:
		return true
	default:
		return false
	}
}

func (p *Parser) parseFunctionOrVariable() (Node, error) {
	identifier := p.currentToken.Value
	p.advance() // Move past the identifier token

	// Check if the next token is an opening parenthesis, indicating a function call
	if p.currentToken.Type == TOKEN_LPAREN {
		return p.parseFunctionCall(identifier)
	}

	// Otherwise, treat the identifier as a variable
	// Assuming you have a VariableNode struct to represent variables
	return p.parseExpression()
}

func (p *Parser) parseOperand() (Node, error) {
	switch p.currentToken.Type {
	case TOKEN_NUMBER:
		value, ok := new(big.Float).SetString(p.currentToken.Value)
		if !ok {
			return nil, fmt.Errorf("invalid number format: %s", p.currentToken.Value)
		}
		node := &ConstantNode{Value: value}
		p.advance() // Move past the number token
		return node, nil
	case TOKEN_LPAREN:
		p.advance()                      // Consume '('
		node, err := p.parseExpression() // Parse the expression inside parentheses
		if err != nil {
			return nil, err
		}
		if !p.expect(TOKEN_RPAREN) {
			return nil, fmt.Errorf("expected closing parenthesis, got %v", p.currentToken)
		}
		return node, nil
	case TOKEN_IDENTIFIER:
		return p.parseFunctionOrVariable()
	default:
		return nil, fmt.Errorf("unexpected token %v in operand", p.currentToken)
	}
}

func (p *Parser) parsePrimary() (Node, error) {
	var node Node
	var err error
	switch {
	case p.currentToken.Type == TOKEN_NUMBER:
		node, err = p.parseNumber()
		if err != nil {
			return nil, err
		}
	case p.currentToken.Type == TOKEN_IDENTIFIER:
		node = &VariableNode{Name: p.currentToken.Value}
		p.advance() // Move past the variable
	case p.currentToken.Type == TOKEN_LPAREN:
		p.advance()                     // Consume '('
		node, err = p.parseExpression() // Parse the sub-expression
		if err != nil {
			return nil, err
		}
		if !p.expect(TOKEN_RPAREN) {
			return nil, fmt.Errorf("expected ')', got %s", p.currentToken.Value)
		}
	default:
		return nil, fmt.Errorf("unexpected token: %s", p.currentToken.Value)
	}

	// Check for a factorial operator after the primary expression
	if p.currentToken.Type == TOKEN_FACTORIAL {
		p.advance() // Consume '!'
		node = &FactorialNode{Operand: node}
	}

	return node, nil
}

func (p *Parser) parseNumber() (Node, error) {
	value, _, err := big.ParseFloat(p.currentToken.Value, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("could not parse number: %s", p.currentToken.Value)
	}
	p.advance() // Move past the number
	return &ConstantNode{Value: value}, nil
}

func (p *Parser) parseFactorial() (Node, error) {
	operand, err := p.parsePrimary() // parsePrimary should parse the primary expression before "!"
	if err != nil {
		return nil, err
	}

	// If the next token is "!", it's a factorial operation
	if p.currentToken.Type == TOKEN_FACTORIAL {
		p.advance() // Consume "!"
		return &FactorialNode{Operand: operand}, nil
	}

	return operand, nil
}

// Parse starts the parsing process and returns the AST.
func (p *Parser) Parse() (Node, error) {
	p.advance() // Initialize parsing by loading the first token
	ast, err := p.parseExpression()
	if err != nil {
		fmt.Println("parse error: AST generation failed")
		return nil, err // Consider returning an error as well
	}
	fmt.Println("AST generated successfully:", ast)
	return ast, nil
}

// InvokeCustomFunctionOrAction tries to execute a function or action by name.
func InvokeCustomFunctionOrAction(name string, precision uint, args ...*big.Float) ([]*big.Float, error) {
	// First, check if it's a mathematical function with a known signature.
	if function, ok := customFuncsMap[name]; ok {
		result, err := function(precision, args...)
		if err != nil {
			return nil, err
		}
		// Wrap single *big.Float result in a slice for a unified interface.
		return []*big.Float{result}, nil
	}

	// Next, check if it's a custom action which we assume returns a slice of *big.Float for this example.
	if action, ok := customActionsMap[name]; ok {
		// Assuming all actions follow a specific signature for this example.
		actionFunc, ok := action.(func(precision uint, args ...*big.Float) ([]*big.Float, error))
		if !ok {
			return nil, fmt.Errorf("action '%s' does not match the expected signature", name)
		}
		return actionFunc(precision, args...)
	}

	return nil, fmt.Errorf("no function or action named '%s' found", name)
}

func ParseSimple(expression string) (float64, error) {
	result, err := parse(expression)
	if err != nil {
		return 0, err
	}

	result64, _ := result.Float64()
	return result64, nil
}

func ParseWithVariables(expression string, variables map[string]float64) (float64, error) {
	substitutedExpression := substituteVariables(expression, variables)
	result, err := parse(substitutedExpression)
	if err != nil {
		return 0, err
	}

	result64, _ := result.Float64()
	return result64, nil
}

// ParseWithVariablesAndReturnString parses the expression and returns its result as a string.
func ParseWithVariablesAndReturnString(expression string, variables map[string]float64) (string, error) {
	substitutedExpression := substituteVariables(expression, variables)
	result, err := parse(substitutedExpression)
	if err != nil {
		return "", err
	}
	// Using String() to get the exact representation without precision loss.
	return result.String(), nil
}

// ParseAndReturnString parses the expression and returns its result as a string.
func ParseAndReturnString(expression string) (string, error) {
	result, err := parse(expression)
	if err != nil {
		return "", err
	}
	// Using String() to get the exact representation without precision loss.
	return result.String(), nil
}

// ParseAndReturnBigDecimal parses the expression and returns its result as a *big.Float.
func ParseAndReturnBigDecimal(expression string) (*big.Float, error) {
	return parse(expression) // Directly return the *big.Float result.
}

// ParseAndReturnBytes parses the expression and returns its result as bytes.
func ParseAndReturnBytes(expression string) ([]byte, error) {
	result, err := parse(expression)
	if err != nil {
		return nil, err
	}
	// Marshalling the *big.Float into bytes.
	return result.GobEncode()
}

// parse is an unexported helper function to handle the common parsing logic.
func parse(expression string) (*big.Float, error) {
	lexer := NewLexer(expression) // Initialize the lexer with the expression.
	parser := NewParser(lexer)    // Create a new parser instance with the lexer.

	// Start parsing by loading the first token.
	parser.advance()

	// Parse the expression to generate the AST. Adjust this to your actual parse method.
	ast, _ := parser.parseExpression() // Assuming this method exists and is implemented correctly.
	if ast == nil {
		fmt.Println("parse error: AST generation failed")
		return nil, fmt.Errorf("parse error: AST generation failed")
	}
	fmt.Println("ast: ", ast)

	// Evaluate the AST and return the result.
	return evaluate(ast)
}

func substituteVariables(expression string, variables map[string]float64) string {
	for varName, varValue := range variables {
		expression = strings.Replace(expression, varName, strconv.FormatFloat(varValue, 'f', -1, 64), -1)
	}
	return expression
}

func evaluate(node Node) (*big.Float, error) {
	fmt.Println("node: ", node)
	if node == nil {
		fmt.Println("node == nil", node)
		return nil, fmt.Errorf("node is nil")
	}
	const epsilon = 1e-10
	switch n := node.(type) {

	case *FactorialNode:
		operandVal, err := evaluate(n.Operand)
		if err != nil {
			return nil, err
		}
		// Convert operandVal to an integer for factorial calculation
		operandInt, _ := operandVal.Int64()
		result := big.NewInt(1) // Initialize result as 1 for multiplication

		for i := int64(1); i <= operandInt; i++ {
			result.Mul(result, big.NewInt(i))
		}

		return new(big.Float).SetInt(result), nil

	case *VariableNode:
		if value, ok := variableContext[n.Name]; ok {
			return value, nil
		} else {
			return nil, fmt.Errorf("variable '%s' not defined", n.Name)
		}

	case *SummationNode:
		sum := big.NewFloat(0).SetPrec(precision)
		startVal, err := evaluate(n.Start)
		if err != nil {
			return nil, err
		}
		var endVal *big.Float
		if _, ok := n.End.(*InfinityNode); !ok {
			endVal, err = evaluate(n.End)
			if err != nil {
				return nil, err
			}
		}

		// Temporarily store the original value of the variable, if it exists
		origVal, origValExists := variableContext[n.Variable]

		// Evaluate the summation
		for i := new(big.Float).SetPrec(precision).Set(startVal); endVal == nil || i.SetPrec(precision).Cmp(endVal) <= 0; i.SetPrec(precision).Add(i, big.NewFloat(1)) {
			// Set the current iteration value in the variable context
			variableContext[n.Variable] = i

			termVal, err := evaluate(n.Expression)
			if err != nil {
				return nil, err
			}
			sum.Add(sum, termVal)

			// Check for convergence if the series is infinite
			if endVal == nil && termVal.SetPrec(precision).Abs(termVal).Cmp(big.NewFloat(epsilon)) < 0 {
				break // Assumed convergence
			}
		}

		// Restore the original variable value, if it was previously defined
		if origValExists {
			variableContext[n.Variable] = origVal
		} else {
			delete(variableContext, n.Variable) // Clean up the temporary variable entry
		}

		return sum, nil

	case *ConstantNode:
		fmt.Println("ConstantNode n: ", n)
		fmt.Println("ConstantNode n.Value: ", n.Value)
		val, _, err := new(big.Float).SetPrec(precision).Parse(n.Value.Text('f', -1), 10)
		if err != nil {
			return nil, err
		}
		return val, nil
	case *FunctionCallNode:
		fmt.Println("FunctionCallNode n: ", n)
		fmt.Println("FunctionCallNode n.Value: ", n.Arguments)
		args := make([]*big.Float, len(n.Arguments))
		for i, argNode := range n.Arguments {
			argVal, err := evaluate(argNode)
			if err != nil {
				return nil, err
			}
			args[i] = argVal
		}
		if function, exists := customFuncsMap[n.Name]; exists {
			return function(precision, args...)
		}
		if action, exists := customActionsMap[n.Name]; exists {
			actionFunc, ok := action.(func(precision uint, args ...*big.Float) (*big.Float, error))
			if !ok {
				return nil, fmt.Errorf("incorrect signature for action '%s'", n.Name)
			}
			return actionFunc(precision, args...)
		}

		return nil, fmt.Errorf("unknown function: %s", n.Name)
	case *UnaryOpNode:
		fmt.Println("UnaryOpNode n: ", n)
		fmt.Println("UnaryOpNode n.Value: ", n.Operand)
		operandVal, err := evaluate(n.Operand)
		if err != nil {
			return nil, err
		}
		if n.Op == "-" { // Negation
			return new(big.Float).SetPrec(precision).Neg(operandVal), nil
		}
		return nil, fmt.Errorf("unsupported unary operation: %s", n.Op)
	case *BinaryOpNode:
		fmt.Println("BinaryOpNode n: ", n)
		fmt.Println("BinaryOpNode n.Left: ", n.Left)
		fmt.Println("BinaryOpNode n.Right: ", n.Right)
		if n.Left == nil || n.Right == nil {
			fmt.Println("BinaryOpNode's left or right child is nil")
			return nil, fmt.Errorf("BinaryOpNode's left or right child is nil")
		}
		leftVal, err := evaluate(n.Left)
		if err != nil {
			return nil, err
		}
		rightVal, err := evaluate(n.Right)
		if err != nil {
			return nil, err
		}

		switch n.Op {
		case "+":
			return new(big.Float).SetPrec(precision).Add(leftVal, rightVal), nil
		case "-":
			return new(big.Float).SetPrec(precision).Sub(leftVal, rightVal), nil
		case "*":
			return new(big.Float).SetPrec(precision).Mul(leftVal, rightVal), nil
		case "/":
			if rightVal.SetPrec(precision).Cmp(new(big.Float).SetFloat64(0)) == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return new(big.Float).SetPrec(precision).Quo(leftVal, rightVal), nil
		case "^":
			// Assuming BinaryExponentiation is correctly implemented elsewhere to handle *big.Float
			return basicMathematicalFunctions.BinaryExponentiation(precision, leftVal, rightVal)
		default:
			return nil, fmt.Errorf("unsupported binary operation: %s", n.Op)
		}
	case *ConditionalNode:
		fmt.Println("ConditionalNode n: ", n)
		fmt.Println("ConditionalNode n.Condition: ", n.Condition)
		conditionVal, err := evaluate(n.Condition)
		if err != nil {
			return nil, err
		}

		// Assuming the condition is a boolean. You'll need a way to interpret the *big.Float as a boolean.
		// This is just a conceptual example. Actual implementation will depend on how you define and evaluate conditions.
		if conditionVal.Cmp(big.NewFloat(1)) == 0 { // True
			return evaluate(n.Consequence)
		} else { // False
			return evaluate(n.Alternative)
		}
	default:
		fmt.Println("default n: ", n)
		return nil, fmt.Errorf("unsupported node type")
	}
}
