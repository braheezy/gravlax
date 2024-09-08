package lox

import "fmt"

func execute(stmt Stmt) *RuntimeError {
	return stmt.Execute()
}

func (p Print) Execute() *RuntimeError {
	value, err := p.expression.Eval()
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", stringify(value))
	return nil
}

func (r Return) Execute() *RuntimeError {
	var value interface{}
	if r.value != nil {
		value, _ = r.value.Eval()
	}

	panic(NewReturn(value))
}

func (e Expression) Execute() *RuntimeError {
	_, err := e.expression.Eval()
	if err != nil {
		return err
	}
	return nil
}
func (f Function) Execute() *RuntimeError {
	fun := LoxFunction{&f}
	interpreter.environment.define(f.name.Lexeme, fun)
	return nil
}

func (i If) Execute() *RuntimeError {
	val, err := i.condition.Eval()
	if isTruthy(val) {
		err = i.thenBranch.Execute()
	} else if i.elseBranch != nil {
		err = i.elseBranch.Execute()
	}
	return err
}

func (v Var) Execute() *RuntimeError {
	var value interface{}
	var err *RuntimeError
	if v.initializer != nil {
		value, err = v.initializer.Eval()
		if err != nil {
			return err
		}
	}

	interpreter.environment.define(v.name.Lexeme, value)
	return nil
}
func (w While) Execute() *RuntimeError {
	for {
		val, _ := w.condition.Eval()
		if !isTruthy(val) {
			break
		}

		err := w.body.Execute()
		if err != nil {
			if err.Message == "break" {
				break
			}
			return err
		}
	}
	return nil
}
func (b Block) Execute() *RuntimeError {
	return executeBlock(b.statements, NewEnvironmentWithEnclosing(interpreter.environment))
}
func executeBlock(statements []Stmt, env *Environment) *RuntimeError {
	previous := interpreter.environment

	defer func() {
		interpreter.environment = previous
	}()

	interpreter.environment = env
	for _, stmt := range statements {
		err := stmt.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (b Break) Execute() *RuntimeError {
	return &RuntimeError{
		Message: "break",
	}
}
