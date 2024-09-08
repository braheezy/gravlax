package lox

type Interpreter struct {
	globals     *Environment
	environment *Environment
}

func NewInterpreter() *Interpreter {
	i := Interpreter{}
	i.globals = NewEnvironment()
	i.environment = i.globals

	i.globals.define("clock", ClockFunction{})

	return &i
}

func (i *Interpreter) interpret(statements []Stmt) {
	for _, statement := range statements {

		err := execute(statement)
		if err != nil {
			handleRuntimeError(err)
		}
	}
}
