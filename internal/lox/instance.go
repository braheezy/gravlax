package lox

import "fmt"

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func (li *LoxInstance) toString() string {
	return li.class.name + " instance"
}

func (li *LoxInstance) get(name Token) (interface{}, *RuntimeError) {
	if value, exists := li.fields[name.Lexeme]; exists {
		return value, nil
	}
	method := li.class.findMethod(name.Lexeme)
	if method != nil {
		return method.bind(li), nil
	}

	return nil, &RuntimeError{name, fmt.Sprintf("Undefined property '%v'.", name.Lexeme)}
}

func (li *LoxInstance) set(name Token, value interface{}) {
	li.fields[name.Lexeme] = value
}
