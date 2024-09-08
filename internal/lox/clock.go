package lox

import "time"

// ClockFunction implements the Callable interface, acting as a native function.
type ClockFunction struct{}

// call method returns the current time in seconds since the epoch.
func (c ClockFunction) call(arguments []interface{}) interface{} {
	// Return the current time in seconds as a floating-point number
	return float64(time.Now().UnixMilli()) / 1000.0
}

// arity method returns 0 because this function expects no arguments.
func (c ClockFunction) arity() int {
	return 0
}

// String method provides a string representation of the function.
func (c ClockFunction) String() string {
	return "<native fn>"
}
