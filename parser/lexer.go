package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// InputToken represents a lexical token.
type InputToken struct {
	Type  int
	Value string
}

func (t *InputToken) isOperatorToken() bool {
	switch t.Type {
	case TOKEN_PLUS, TOKEN_MINUS, TOKEN_ASTERISK, TOKEN_SLASH, TOKEN_CARET:
		return true
	default:
		return false
	}
}

// Lexer breaks the input into tokens.
type Lexer struct {
	input         string
	start         int
	pos           int
	length        int
	lastTokenType int
}

// NewLexer creates a new instance of Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, length: len(input)}
}

func (l *Lexer) prevTokenIsNumeric() bool {
	return l.lastTokenType == TOKEN_NUMBER
}

func (l *Lexer) readNumber() InputToken {
	startPosition := l.pos
	hasDecimal := false
	for l.pos < l.length && (unicode.IsDigit(rune(l.input[l.pos])) || (l.input[l.pos] == '.' && !hasDecimal)) {
		if l.input[l.pos] == '.' {
			hasDecimal = true
		}
		l.pos++
	}
	// No additional advance needed here; l.pos is already at the next character
	return InputToken{Type: TOKEN_NUMBER, Value: l.input[startPosition:l.pos]}
}

func (l *Lexer) readIdentifier() InputToken {
	startPosition := l.pos
	for l.pos < l.length && (isLetter(l.input[l.pos]) || unicode.IsDigit(rune(l.input[l.pos]))) {
		l.pos++
	}
	// No additional advance needed here; l.pos is already at the next character
	return InputToken{Type: TOKEN_IDENTIFIER, Value: l.input[startPosition:l.pos]}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() InputToken {

	fmt.Printf("Starting NextToken at pos: %d\n", l.pos)
	for l.pos < l.length {
		fmt.Printf("Current char token: '%c', pos: %d\n", l.input[l.pos], l.pos)
		switch {
		case unicode.IsSpace(rune(l.input[l.pos])):
			fmt.Println("Found whitespace")
			l.ignore()
		case strings.HasPrefix(l.input[l.pos:], "sum of"):
			fmt.Println("Found 'sum of'")
			return l.consumeToken(TOKEN_SUM_OF, "sum of")
		case strings.HasPrefix(l.input[l.pos:], "from"):
			fmt.Println("Found 'from'")
			return l.consumeToken(TOKEN_FROM, "from")
		case strings.HasPrefix(l.input[l.pos:], "to"):
			fmt.Println("Found 'to'")
			return l.consumeToken(TOKEN_TO, "to")
		case strings.HasPrefix(l.input[l.pos:], "infinity"):
			fmt.Println("Found 'infinity'")
			return l.consumeToken(TOKEN_INFINITY, "infinity")
		case l.input[l.pos] == '(':
			fmt.Println("Found '('")
			return l.consumeToken(TOKEN_LPAREN, "(")
		case l.input[l.pos] == ')':
			fmt.Println("Found ')'")
			return l.consumeToken(TOKEN_RPAREN, ")")
		case l.input[l.pos] == '+':
			fmt.Println("Found '+'")
			return l.consumeToken(TOKEN_PLUS, "+")
		case l.input[l.pos] == '-':
			fmt.Println("Found '-'")
			return l.consumeToken(TOKEN_MINUS, "-")
		case l.input[l.pos] == '*':
			fmt.Printf("Found '*', pos: %d\n", l.pos)
			return l.consumeToken(TOKEN_ASTERISK, "*")
		case l.input[l.pos] == '/':
			fmt.Println("Found '/'")
			return l.consumeToken(TOKEN_SLASH, "/")
		case l.input[l.pos] == '^':
			fmt.Println("Found '^'")
			return l.consumeToken(TOKEN_CARET, "^")
		case l.input[l.pos] == '=':
			fmt.Println("Found '='")
			return l.consumeToken(TOKEN_EQUALS, "=")
		case l.input[l.pos] == '!':
			fmt.Println("Found '!'")
			return l.consumeToken(TOKEN_FACTORIAL, "!")
		case unicode.IsDigit(rune(l.input[l.pos])) || l.input[l.pos] == '.':
			fmt.Println("Found number or '.'")
			l.lastTokenType = TOKEN_NUMBER
			return l.readNumber()
		case isLetter(l.input[l.pos]):
			fmt.Println("Found identifier")
			if l.prevTokenIsNumeric() {
				// If the previous token was numeric and the current character can start an identifier,
				// it indicates a new token starts here. We should not handle it as a part of the number.
				l.handlePotentialOperatorOrNewToken()
			}
			token := l.readIdentifier()
			l.lastTokenType = token.Type
			return token
		default:
			fmt.Printf("Skipping unrecognized character: %c\n", l.input[l.pos])
			l.advance(1) // Skip characters that don't match any case
		}
	}
	//fmt.Println("Reached EOF")
	return InputToken{Type: TOKEN_EOF, Value: ""}
}

func (l *Lexer) handlePotentialOperatorOrNewToken() {
	// This function would be called only if prevTokenIsNumeric() returns true
	// and the current character is an identifier or '(' that directly follows a number.
	// Since the current implementation directly reads and advances the lexer,
	// the actual handling logic might need to be integrated into your parsing strategy,
	// especially if you want to insert an implicit multiplication token.

	// Placeholder for logic to handle the situation.
	// This might include setting a flag, adjusting the lexer's position,
	// or directly injecting tokens into the parsing process.

	// Example: Resetting the start position of the lexer to re-evaluate the current character
	// as the beginning of a new token. This is a simplistic approach and might not cover all cases.
	// You may need a more sophisticated state management or token injection strategy.
	l.start = l.pos
}

func (l *Lexer) consumeToken(tokenType int, value string) InputToken {
	fmt.Printf("Consuming token: %s\n", value)
	token := InputToken{Type: tokenType, Value: value}
	l.advance(len(value)) // Ensure this correctly advances based on the token's length
	return token
}

func (l *Lexer) ignore() {
	l.pos++
	l.start = l.pos
}

func (l *Lexer) advance(n int) {
	l.pos += n
}
