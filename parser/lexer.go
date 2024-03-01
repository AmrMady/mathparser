package parser

import (
	"fmt"
	"go/token"
	"strings"
	"unicode"
)

// InputToken represents a lexical token.
type InputToken struct {
	Type  int
	Value string
}

// Lexer breaks the input into tokens.
type Lexer struct {
	input  string
	start  int
	pos    int
	length int
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

// NewLexer creates a new instance of Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, length: len(input)}
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
	if startPosition == l.pos { // No digits were read
		return InputToken{Type: TOKEN_ERROR, Value: "Invalid number"}
	}
	return InputToken{Type: TOKEN_NUMBER, Value: l.input[startPosition:l.pos]}
}

func (l *Lexer) readIdentifier() InputToken {
	startPosition := l.pos
	for l.pos < l.length && (isLetter(l.input[l.pos]) || unicode.IsDigit(rune(l.input[l.pos]))) {
		l.advance(1)
	}
	return InputToken{Type: TOKEN_IDENTIFIER, Value: l.input[startPosition:l.pos]}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() InputToken {
	if l.pos >= l.length {
		return InputToken{Type: TOKEN_EOF, Value: ""}
	}
	// Before recognizing "to", explicitly print the current position and character
	fmt.Printf("About to check 'to' at pos %d: '%s'\n", l.pos, l.input[l.pos:])
	for l.pos < l.length {
		switch {
		case unicode.IsSpace(rune(l.input[l.pos])):
			l.ignore()
		case strings.HasPrefix(l.input[l.pos:], "sum of"):
			l.advance(len("sum of"))
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_SUM_OF, Value: "sum of"}) // Debug print
			return InputToken{Type: TOKEN_SUM_OF, Value: "sum of"}
		case strings.HasPrefix(l.input[l.pos:], "from"):
			l.advance(len("from"))
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_FROM, Value: "from"}) // Debug print
			return InputToken{Type: TOKEN_FROM, Value: "from"}
		case strings.HasPrefix(l.input[l.pos:], "to"):
			l.advance(len("to"))
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_TO, Value: "to"}) // Debug print
			return InputToken{Type: TOKEN_TO, Value: "to"}
		case strings.HasPrefix(l.input[l.pos:], "infinity"):
			l.advance(len("infinity"))
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_INFINITY, Value: "infinity"}) // Debug print
			return InputToken{Type: TOKEN_INFINITY, Value: "infinity"}
		case l.input[l.pos] == '(':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_LPAREN, Value: "("}) // Debug print
			return InputToken{Type: TOKEN_LPAREN, Value: "("}
		case l.input[l.pos] == ')':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_RPAREN, Value: ")"}) // Debug print
			return InputToken{Type: TOKEN_RPAREN, Value: ")"}
		case l.input[l.pos] == '+':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_PLUS, Value: "+"}) // Debug print
			return InputToken{Type: TOKEN_PLUS, Value: "+"}
		case l.input[l.pos] == '-':
			// Special handling for unary minus might be needed, depending on context
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_MINUS, Value: "-"}) // Debug print
			return InputToken{Type: TOKEN_MINUS, Value: "-"}
		case l.input[l.pos] == '*':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_ASTERISK, Value: "*"}) // Debug print
			return InputToken{Type: TOKEN_ASTERISK, Value: "*"}
		case l.input[l.pos] == '/':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_SLASH, Value: "/"}) // Debug print
			return InputToken{Type: TOKEN_SLASH, Value: "/"}
		case l.input[l.pos] == '^':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_CARET, Value: "^"}) // Debug print
			return InputToken{Type: TOKEN_CARET, Value: "^"}
		case l.input[l.pos] == '=':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_EQUALS, Value: "="}) // Debug print
			return InputToken{Type: TOKEN_EQUALS, Value: "="}
		case l.input[l.pos] == '!':
			l.advance(1)
			fmt.Printf("Generated token: %v\n", InputToken{Type: TOKEN_FACTORIAL, Value: "!"}) // Debug print
			return InputToken{Type: TOKEN_FACTORIAL, Value: "!"}
		case unicode.IsDigit(rune(l.input[l.pos])): // Start with a digit
			return l.readNumber() // This handles both integers and floats.
		case l.input[l.pos] == '.':
			// Ensure there's a digit following the '.', indicating a float.
			if l.pos+1 < l.length && unicode.IsDigit(rune(l.input[l.pos+1])) {
				return l.readNumber()
			}
		// Handling identifiers (e.g., function names, variables)
		case isLetter(l.input[l.pos]):
			fmt.Printf("Generated token: %v\n", l.readIdentifier()) // Debug print
			return l.readIdentifier()
		// Add cases for other tokens...
		default:
			l.advance(1) // Skip characters that don't match any case
		}
	}
	// After recognizing "to", print the next character
	fmt.Printf("Next char after 'to': '%s'\n", l.input[l.pos:])
	fmt.Printf("Generated token: %v\n", InputToken{Type: int(token.EOF), Value: ""}) // Debug print
	return InputToken{Type: int(token.EOF), Value: ""}
}

func (l *Lexer) ignore() {
	l.pos++
	l.start = l.pos
}

func (l *Lexer) advance(n int) {
	l.pos += n
}
