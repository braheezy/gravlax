package lox

import (
	"errors"
	"fmt"
)

type Parser struct {
	Tokens    []Token
	current   int
	loopDepth int
}

func (p *Parser) Parse() ([]Stmt, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			// Handle panic if it's a ParseError
			if parseErr, ok := r.(ParseError); ok {
				err = parseErr
			} else {
				panic(r) // If it's not a ParseError, re-panic
			}
		}
	}()

	var statements []Stmt
	for !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, dec)
	}
	return statements, err
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() (stmt Stmt, err ParseError) {
	defer func() {
		if r := recover(); r != nil {
			// Handle panic if it's a ParseError
			if parseError, ok := r.(ParseError); ok {
				p.synchronize()
				err = parseError
			} else {
				panic(r)
			}
		}
	}()

	if p.match(FUN) {
		return p.function("function"), nil
	}
	if p.match(VAR) {
		stmt = p.varDeclaration()
	} else {
		stmt = p.statement()
	}
	return stmt, err
}
func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(BREAK) {
		return p.breakStatement()
	}
	if p.match(LEFT_BRACE) {
		return Block{statements: p.block()}
	}
	return p.expressionStatement()
}
func (p *Parser) forStatement() Stmt {
	p.loopDepth++
	defer func() { p.loopDepth-- }()

	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = Block{[]Stmt{body, Expression{increment}}}
	}
	if condition == nil {
		condition = Literal{true}
	}
	body = While{condition, body}

	if initializer != nil {
		body = Block{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return If{condition, thenBranch, elseBranch}

}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Print{expression: value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return Return{keyword, value}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return Var{name: name, initializer: initializer}
}

func (p *Parser) whileStatement() Stmt {
	p.loopDepth++
	defer func() {}()

	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")

	p.loopDepth++
	body := p.statement()
	p.loopDepth--

	return While{condition, body}
}

func (p *Parser) breakStatement() Stmt {
	if p.loopDepth == 0 {
		p.error(p.previous(), "Cannot use 'break' outside of a loop.")
	}

	p.consume(SEMICOLON, "Expect ';' after 'break'.")
	return Break{}
}

func (p *Parser) expressionStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return Expression{expression: value}
}
func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %v name.", kind))

	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %v name.", kind))
	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %v body.", kind))
	body := p.block()

	return Function{&name, parameters, body}
}
func (p *Parser) anonFunction() Expr {
	p.consume(LEFT_PAREN, "Expect '(' before anonymous function parameters.")
	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, "Expect '{' before anonymous function body.")
	body := p.block()

	// Return the anonymous function as an expression
	return AnonFunction{params: parameters, body: body}
}
func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		dec, _ := p.declaration()
		statements = append(statements, dec)
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if _, ok := expr.(Variable); ok {
			name := expr.(Variable).name
			return Assign{name: name, value: value}
		}

		reportTokenError(equals, "Invalid assignment target.")
	}
	return expr
}
func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = Logical{expr, operator, right}
	}
	return expr
}
func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = Logical{expr, operator, right}
	}
	return expr
}
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = Binary{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return Unary{operator, right}
	}
	return p.call()
}
func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(COMMA) {
			arguments = append(arguments, p.expression())
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return Call{callee, paren, arguments}
}
func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}
func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return Literal{false}
	}
	if p.match(TRUE) {
		return Literal{true}
	}
	if p.match(NIL) {
		return Literal{nil}
	}
	if p.match(NUMBER, STRING) {
		return Literal{p.previous().Literal}
	}
	if p.match(FUN) {
		// Directly parse the anonymous function as an expression
		return p.anonFunction()
	}
	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return Grouping{expr}
	}
	if p.match(IDENTIFIER) {
		return Variable{p.previous()}
	}

	// on a token that can't start an expression
	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(tokenType TokenType, message string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.current-1]
}

type ParseError error

func (p *Parser) error(token Token, message string) error {
	reportTokenError(token, message)
	return ParseError(errors.New(message))
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
