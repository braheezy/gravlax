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

func interpret(statements []Stmt) {
	for _, statement := range statements {

		err := execute(statement)
		if err != nil {
			handleRuntimeError(err)
		}
	}
}

func (l Literal) Eval() (interface{}, *RuntimeError) {
	return l.value, nil
}

func (l Logical) Eval() (interface{}, *RuntimeError) {
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

func (g Grouping) Eval() (interface{}, *RuntimeError) {
	return g.expression.Eval()
}

func (b Binary) Eval() (interface{}, *RuntimeError) {
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

func (u Unary) Eval() (interface{}, *RuntimeError) {
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

func (v Variable) Eval() (interface{}, *RuntimeError) {
	return environment.get(v.name)
}

func (a Assign) Eval() (interface{}, *RuntimeError) {
	value, err := a.value.Eval()
	if err != nil {
		return nil, err
	}

	environment.assign(a.name, value)
	return value, nil
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

func stringify(value interface{}) string {
	switch v := value.(type) {
	case float64:
		// Convert the float32 to a string
		text := fmt.Sprintf("%f", v)
		// If it ends with ".0", remove it
		text = strings.TrimSuffix(text, ".000000")
		return text
	default:
		// Fallback for other types, using fmt.Sprintf to handle them
		return fmt.Sprintf("%v", value)
	}
}

func handleRuntimeError(err *RuntimeError) {
	fmt.Fprintf(os.Stderr, "%s\n[line %d]\n", err.Error(), err.Token.Line)
}
