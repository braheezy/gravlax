package lox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t *Token) String() string {
	return fmt.Sprintf("%v %v %v", tokenNames[t.Type], t.Lexeme, t.Literal)
}
