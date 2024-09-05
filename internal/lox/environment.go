package lox

import "fmt"

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

var environment = NewEnvironment()

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]interface{})}
}
func NewEnvironmentWithEnclosing(enclosing *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}), enclosing: enclosing}
}
func (e *Environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (interface{}, *RuntimeError) {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, &RuntimeError{name, fmt.Sprintf("Undefined variable '%v'.", name.Lexeme)}
}

func (e *Environment) assign(name Token, value interface{}) *RuntimeError {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}

	if e.enclosing != nil {
		e.enclosing.assign(name, value)
		return nil
	}

	return &RuntimeError{name, fmt.Sprintf("Undefined variable %v", name.Lexeme)}
}
