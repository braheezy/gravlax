package lox

import (
	"fmt"
	"os"
	"strings"
)

// Assuming RuntimeError is a custom error type
type RuntimeError struct {
	Token   Token
	Message string
}

// Implement the Error() method to satisfy the error interface
func (e *RuntimeError) Error() string {
	return e.Message
}

func (l *Literal) Eval() (interface{}, *RuntimeError) {
	return l.value, nil
}
func (l *Logical) Eval() (interface{}, *RuntimeError) {
	left, _ := l.left.Eval()

	if l.operator.Type == OR {
		if isTruthy(left) {
			return left, nil
		} else {
			if !isTruthy(left) {
				return left, nil
			}
		}
	}
	return l.right.Eval()
}
func (s *Set) Eval() (interface{}, *RuntimeError) {
	object, err := s.object.Eval()
	if err != nil {
		return nil, err
	}

	if _, ok := object.(*LoxInstance); !ok {
		return nil, &RuntimeError{s.name, "Only instances have fields."}
	}

	value, err := s.value.Eval()
	if err != nil {
		return nil, err
	}
	object.(*LoxInstance).set(s.name, value)
	return value, nil
}
func (t *This) Eval() (interface{}, *RuntimeError) {
	return lookupVariable(t.keyword, t)
}
func (g *Grouping) Eval() (interface{}, *RuntimeError) {
	return g.expression.Eval()
}
func (b *Binary) Eval() (interface{}, *RuntimeError) {
	left, err := b.left.Eval()
	if err != nil {
		return nil, err
	}
	right, err := b.right.Eval()
	if err != nil {
		return nil, err
	}

	leftNumber, leftOk := left.(float64)

	rightNumber, rightOK := right.(float64)

	switch b.operator.Type {
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case GREATER:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber > rightNumber, nil
	case GREATER_EQUAL:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber >= rightNumber, nil
	case LESS:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber < rightNumber, nil
	case LESS_EQUAL:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber <= rightNumber, nil
	case MINUS:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber - rightNumber, nil
	case SLASH:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber / rightNumber, nil
	case STAR:
		if !leftOk || !rightOK {
			return nil, &RuntimeError{b.operator, "operands must be numbers"}
		}
		return leftNumber * rightNumber, nil
	case PLUS:
		if leftOk && rightOK {
			return leftNumber + rightNumber, nil
		}
		leftString, leftOk := left.(string)
		rightString, rightOK := right.(string)

		if leftOk && rightOK {
			return leftString + rightString, nil
		}
		return nil, &RuntimeError{b.operator, "operands must be two numbers or two strings"}
	}

	return nil, nil
}
func (c *Call) Eval() (interface{}, *RuntimeError) {
	callee, _ := c.callee.Eval()

	var arguments []interface{}
	for _, arg := range c.arguments {
		val, _ := arg.Eval()
		arguments = append(arguments, val)
	}

	function, ok := callee.(Callable)
	if !ok {
		return nil, &RuntimeError{
			// Assuming c.paren is the token that represents the call
			Token:   c.paren,
			Message: "Can only call functions and classes.",
		}
	}

	if len(arguments) != function.arity() {
		return nil, &RuntimeError{
			// Assuming c.paren is the token that represents the call
			Token:   c.paren,
			Message: fmt.Sprintf("Expected %v arguments but got %v.", function.arity(), len(arguments)),
		}
	}
	return function.call(arguments), nil
}
func (g *Get) Eval() (interface{}, *RuntimeError) {
	object, err := g.object.Eval()
	if err != nil {
		return nil, err
	}
	if inst, ok := object.(*LoxInstance); ok {
		return inst.get(g.name)
	}
	return nil, &RuntimeError{g.name, "Only instances have properties."}
}
func (u *Unary) Eval() (interface{}, *RuntimeError) {
	right, err := u.right.Eval()
	if err != nil {
		return nil, err
	}
	switch u.operator.Type {
	case MINUS:
		if number, ok := right.(float64); ok {
			return -number, nil
		} else {
			return nil, &RuntimeError{u.operator, "operand must be a number"}
		}
	case BANG:
		return !isTruthy(right), nil
	}
	return nil, nil
}
func (v *Variable) Eval() (interface{}, *RuntimeError) {
	return lookupVariable(v.name, v)
}
func lookupVariable(name Token, expr Expr) (interface{}, *RuntimeError) {
	distance, exists := interpreter.locals[expr]
	if exists {
		return interpreter.environment.getAt(distance, name.Lexeme)
	} else {
		return interpreter.globals.get(name)
	}
}
func (a *Assign) Eval() (interface{}, *RuntimeError) {
	value, err := a.value.Eval()
	if err != nil {
		return nil, err
	}

	distance, exists := interpreter.locals[a]
	if exists {
		interpreter.environment.assignAt(distance, a.name, value)
	} else {
		interpreter.globals.assign(a.name, value)
	}
	return value, nil
}
func (af *AnonFunction) Eval() (interface{}, *RuntimeError) {
	return &LoxFunction{
		declaration: &Function{
			name:   nil, // Anonymous functions have no name
			params: af.params,
			body:   af.body,
		},
	}, nil
}
func isTruthy(e interface{}) bool {
	if e == nil {
		return false
	}
	if value, ok := e.(bool); ok {
		return value
	}
	return true
}

func isEqual(a interface{}, b interface{}) bool {
	return a == b
}

type Stringer interface {
	toString() string
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case float64:
		// Convert the float32 to a string
		text := fmt.Sprintf("%f", v)
		// If it ends with ".0", remove it
		text = strings.TrimSuffix(text, ".000000")
		return text
	case Stringer:
		return v.toString()
	default:
		// Fallback for other types, using fmt.Sprintf to handle them
		return fmt.Sprintf("%v", value)
	}
}

func handleRuntimeError(err *RuntimeError) {
	fmt.Fprintf(os.Stderr, "[line %d]{%v} %s\n", err.Token.Line, err.Token.Lexeme, err.Error())
}
