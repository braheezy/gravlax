package lox

import "fmt"

type Callable interface {
	call(*Interpreter, []interface{}) interface{}
	arity() int
	toString() string
}

type LoxFunction struct {
	declaration *Function
}

func (lf LoxFunction) call(i *Interpreter, arguments []interface{}) (out interface{}) {
	environment := NewEnvironmentWithEnclosing(i.globals)

	for i := 0; i < len(lf.declaration.params); i++ {
		environment.define(lf.declaration.params[i].Lexeme, arguments[i])
	}

	// Try to execute the function block and catch the return value if it occurs.
	defer func() {
		if r := recover(); r != nil {
			if returnValue, ok := r.(*Ret); ok {
				// Unwind the call with the return value.
				out = returnValue.value
			} else {
				panic(r) // Re-panic if it's not the expected return value error.
			}
		}
	}()

	executeBlock(lf.declaration.body, environment)
	return nil
}

func (lf LoxFunction) arity() int {
	return len(lf.declaration.params)
}

func (lf LoxFunction) toString() string {
	return fmt.Sprintf("<fn %v>", lf.declaration.name.Lexeme)
}

// Return is a custom error type used to signal a function return.
type Ret struct {
	value interface{}
}

// Error method satisfies the error interface but is not used since it's meant for control flow.
func (r *Ret) Error() string {
	return "<return>"
}

// NewReturn creates a new Return instance with the given value.
func NewReturn(value interface{}) *Ret {
	return &Ret{value: value}
}
