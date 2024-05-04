package main

import (
	"log"
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
	source         string
	tokens         []Token
	start          int
	current        int
	line           int
	inBlockComment bool
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{Type: EOF, Line: s.line})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	if s.inBlockComment {
		// Continue skipping characters until the end of the block comment
		for !s.isAtEnd() {
			if s.peek() == '*' && s.peekNext() == '/' {
				s.advance()              // Advance to '*'
				s.advance()              // Advance to '/'
				s.inBlockComment = false // End block comment
				return
			}
			s.advance()
		}
		return // if EOF reached while still in a block comment
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
			s.inBlockComment = true
			// A block comment goes until another '*/' is encountered
			for !s.isAtEnd() && s.inBlockComment {
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance() // Advance to '*'
					s.advance() // Advance to '/'
					s.inBlockComment = false
				} else {
					s.advance()
				}
			}
			// for {
			// 	if s.peek() == '*' && s.peekNext() == '/' {
			// 		// Consume the closing '*/'
			// 		s.advance()
			// 		s.advance()
			// 		break
			// 	} else if s.isAtEnd() {
			// 		reportError(s.line, "Unterminated block comment")
			// 	} else {
			// 		if s.peek() == '\n' {
			// 			s.advance()
			// 		}
			// 		s.advance()
			// 	}
			// }
			// for s.peek() != '*' && s.peekNext() != '/' && !s.isAtEnd() {
			// 	if s.peek() == '\n' {
			// 		continue
			// 	} else {
			// 		s.advance()
			// 	}
			// }
			// Consume the last '*/'
			// s.advance()
			// s.advance()
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.handleString()
	default:
		if isDigit(c) {
			s.handleNumber()
		} else if isAlpha(c) {
			s.handleIdentifier()
		} else {
			reportError(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) advance() rune {
	char := s.source[s.current]
	s.current++
	return rune(char)
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) handleString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		reportError(s.line, "Unterminated string.")
		return
	}

	// The closing "
	s.advance()

	// Trim surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, value)
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

	value, err := strconv.Atoi(s.source[s.start:s.current])
	if err != nil {
		log.Fatal(err)
	}
	s.addToken(NUMBER, value)
}

func (s *Scanner) handleIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
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
