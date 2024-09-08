package lox

var globals = NewEnvironment()
var environment = globals

func init() {
	globals.define("clock", ClockFunction{})
}

func interpret(statements []Stmt) {
	for _, statement := range statements {

		err := execute(statement)
		if err != nil {
			handleRuntimeError(err)
		}
	}
}
