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

func (i If) Execute() (interface{}, *RuntimeError) {
	val, err := i.condition.Eval()
	if isTruthy(val) {
		_, err = i.thenBranch.Execute()
	} else if i.elseBranch != nil {
		_, err = i.elseBranch.Execute()
	}
	return nil, err
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
func (w While) Execute() (interface{}, *RuntimeError) {
	for {
		val, _ := w.condition.Eval()
		if !isTruthy(val) {
			break
		}

		_, err := w.body.Execute()
		if err != nil {
			if err.Message == "break" {
				break
			}
			return nil, err
		}
	}
	return nil, nil
}
func (b Block) Execute() (interface{}, *RuntimeError) {
	return nil, executeBlock(b.statements, NewEnvironmentWithEnclosing(environment))
}
func executeBlock(statements []Stmt, env *Environment) *RuntimeError {
	previous := environment

	environment = env
	for _, stmt := range statements {
		_, err := stmt.Execute()
		if err != nil {
			environment = previous
			return err
		}
	}

	environment = previous
	return nil
}

func (b Break) Execute() (interface{}, *RuntimeError) {
	return nil, &RuntimeError{
		Message: "break",
	}
}
