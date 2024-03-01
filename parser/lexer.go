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

// Lexer breaks the input into tokens.
type Lexer struct {
	input  string
	start  int
	pos    int
	length int
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
	fmt.Printf("Starting NextToken at pos: %d\n", l.pos)
	for l.pos < l.length {
		fmt.Printf("Current char: '%c'\n", l.input[l.pos])
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
			fmt.Println("Found '*'")
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
			return l.readNumber()
		case isLetter(l.input[l.pos]):
			fmt.Println("Found identifier")
			return l.readIdentifier()
		default:
			fmt.Printf("Skipping unrecognized character: %c\n", l.input[l.pos])
			l.advance(1) // Skip characters that don't match any case
		}
	}
	fmt.Println("Reached EOF")
	return InputToken{Type: TOKEN_EOF, Value: ""}
}

func (l *Lexer) consumeToken(tokenType int, value string) InputToken {
	fmt.Printf("Consuming token: %s\n", value)
	l.advance(len(value))
	return InputToken{Type: tokenType, Value: value}
}

func (l *Lexer) ignore() {
	l.pos++
	l.start = l.pos
}

func (l *Lexer) advance(n int) {
	l.pos += n
}
