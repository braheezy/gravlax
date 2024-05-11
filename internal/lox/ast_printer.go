package lox

import (
	"fmt"
	"strings"
)

func PrintAST(expression Expr) {
	println(expression.(Stringify).ToString())
}

type Stringify interface {
	ToString() string
}

func parenthesize(label string, exprs ...Expr) string {
	var sb strings.Builder

	sb.WriteRune('(')
	sb.WriteString(label)
	for _, expr := range exprs {
		sb.WriteRune(' ')
		sb.WriteString(expr.(Stringify).ToString())
	}
	sb.WriteRune(')')

	return sb.String()
}

func (l Literal) ToString() string {
	if l.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", l.value)
}

func (b Binary) ToString() string {
	return parenthesize(b.operator.Lexeme, b.left, b.right)
}

func (u Unary) ToString() string {
	return parenthesize(u.operator.Lexeme, u.right)
}

func (g Grouping) ToString() string {
	return parenthesize("group", g.expression)
}
