package lox

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
	hadError    bool
}

var interpreter = NewInterpreter()

func NewInterpreter() *Interpreter {
	i := Interpreter{}
	i.globals = NewEnvironment()
	i.environment = i.globals
	i.locals = make(map[Expr]int)

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

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}
