package lox

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	Source         string
	Tokens         []Token
	start          int
	Current        int
	Line           int
	InBlockComment bool
}

func (s *Scanner) ScanTokens() error {
	var err error
	var hadError bool
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.Current
		err = s.scanToken()
		if err != nil {
			hadError = true
		}
	}

	s.Tokens = append(s.Tokens, Token{Type: EOF, Line: s.Line})

	if hadError {
		return errors.New("")
	}
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= len(s.Source)
}

func (s *Scanner) scanToken() error {
	if s.InBlockComment {
		// Continue skipping characters until the end of the block comment
		for !s.isAtEnd() {
			if s.peek() == '*' && s.peekNext() == '/' {
				s.advance()              // Advance to '*'
				s.advance()              // Advance to '/'
				s.InBlockComment = false // End block comment
				return nil
			}
			s.advance()
		}
		return nil
	}

	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes to the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.InBlockComment = true
			// A block comment goes until another '*/' is encountered
			for !s.isAtEnd() && s.InBlockComment {
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance() // Advance to '*'
					s.advance() // Advance to '/'
					s.InBlockComment = false
				} else {
					s.advance()
				}
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.Line++
	case '"':
		return s.handleString()
	default:
		if isDigit(c) {
			s.handleNumber()
		} else if isAlpha(c) {
			s.handleIdentifier()
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error%s:%s\n", s.Line, "", "Unexpected character.")
			return errors.New("")
		}
	}
	return nil
}

func (s *Scanner) advance() rune {
	char := s.Source[s.Current]
	s.Current++
	return rune(char)
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.Source[s.start:s.Current]
	s.Tokens = append(s.Tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.Line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.Source[s.Current]) != expected {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.Source[s.Current])
}

func (s *Scanner) peekNext() rune {
	if s.Current+1 >= len(s.Source) {
		return rune(0)
	}
	return rune(s.Source[s.Current+1])
}

func (s *Scanner) handleString() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		fmt.Fprintf(os.Stderr, "[line %d] Error%s:%s\n", s.Line, "", "Unterminated string.")
		return errors.New("")
	}

	// The closing "
	s.advance()

	// Trim surrounding quotes
	value := s.Source[s.start+1 : s.Current-1]
	s.addToken(STRING, value)
	return nil
}

func (s *Scanner) handleNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the '.'
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.Atoi(s.Source[s.start:s.Current])
	if err != nil {
		log.Fatal(err)
	}
	s.addToken(NUMBER, value)
}

func (s *Scanner) handleIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.Source[s.start:s.Current]
	if value, ok := keywords[text]; ok {
		s.addToken(value, nil)
	} else {
		s.addToken(IDENTIFIER, nil)
	}

}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
