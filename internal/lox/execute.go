package lox

import "fmt"

func execute(stmt Stmt) *RuntimeError {
	_, err := stmt.Execute()
	return err
}

func (p Print) Execute() (interface{}, *RuntimeError) {
	value, err := p.expression.Eval()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", stringify(value))
	return nil, nil
}

func (e Expression) Execute() (interface{}, *RuntimeError) {
	_, err := e.expression.Eval()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (v Var) Execute() (interface{}, *RuntimeError) {
	var value interface{}
	var err *RuntimeError
	if v.initializer != nil {
		value, err = v.initializer.Eval()
		if err != nil {
			return nil, err
		}
	}

	environment.define(v.name.Lexeme, value)
	return nil, nil
}
func (b Block) Execute() (interface{}, *RuntimeError) {
	executeBlock(b.statements, NewEnvironmentWithEnclosing(environment))
	return nil, nil
}
func executeBlock(statements []Stmt, env Environment) {
	previous := environment

	environment = env
	for _, stmt := range statements {
		_, err := stmt.Execute()
		if err != nil {
			environment = previous
			return
		}
	}

	environment = previous
}