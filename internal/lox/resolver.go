package lox

type FunctionType int

const (
	NoFunct FunctionType = iota
	Funct
)

type LoopType int

const (
	NoLoop LoopType = iota
	Loop
)

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
	currentLoop     LoopType
}

type Resolvable interface {
	Resolve(r *Resolver)
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     i,
		scopes:          make([]map[string]bool, 0),
		currentFunction: NoFunct,
	}
}

func (b Block) Resolve(r *Resolver) {
	r.beginScope()
	r.resolveStatements(b.statements)
	r.endScope()
}

func (e Expression) Resolve(r *Resolver) {
	e.expression.(Resolvable).Resolve(r)
}

func (v Var) Resolve(r *Resolver) {
	r.declare(v.name)
	if v.initializer != nil {
		v.initializer.(Resolvable).Resolve(r)
	}
	r.define(v.name)
}
func (w While) Resolve(r *Resolver) {
	enclosingLoop := r.currentLoop
	r.currentLoop = Loop

	w.condition.(Resolvable).Resolve(r)
	w.body.(Resolvable).Resolve(r)

	r.currentLoop = enclosingLoop
}
func (a Assign) Resolve(r *Resolver) {
	a.value.(Resolvable).Resolve(r)
	r.resolveLocal(a, a.name)
}
func (b Binary) Resolve(r *Resolver) {
	b.left.(Resolvable).Resolve(r)
	b.right.(Resolvable).Resolve(r)
}
func (c Call) Resolve(r *Resolver) {
	c.callee.(Resolvable).Resolve(r)
	for _, arg := range c.arguments {
		arg.(Resolvable).Resolve(r)
	}
}
func (g Grouping) Resolve(r *Resolver) {
	g.expression.(Resolvable).Resolve(r)
}
func (l Literal) Resolve(r *Resolver) {

}
func (l Logical) Resolve(r *Resolver) {
	l.left.(Resolvable).Resolve(r)
	l.right.(Resolvable).Resolve(r)
}
func (u Unary) Resolve(r *Resolver) {
	u.right.(Resolvable).Resolve(r)
}
func (f Function) Resolve(r *Resolver) {
	r.declare(*f.name)
	r.define(*f.name)

	r.resolveFunction(f, Funct)
}
func (i If) Resolve(r *Resolver) {
	i.condition.(Resolvable).Resolve(r)
	i.thenBranch.(Resolvable).Resolve(r)
	if i.elseBranch != nil {
		i.elseBranch.(Resolvable).Resolve(r)
	}
}
func (p Print) Resolve(r *Resolver) {
	p.expression.(Resolvable).Resolve(r)
}
func (af AnonFunction) Resolve(r *Resolver) {
	r.beginScope()

	// Declare and define each parameter within the function's scope
	for _, param := range af.params {
		r.declare(param)
		r.define(param)
	}

	// Resolve the function body statements
	r.resolveStatements(af.body)

	// End the scope after resolving the function body
	r.endScope()
}
func (re Return) Resolve(r *Resolver) {
	if r.currentFunction == NoFunct {
		loxError(re.keyword, "Can't return from top-level code.")
	}

	if re.value != nil {
		re.value.(Resolvable).Resolve(r)
	}
}
func (v Variable) Resolve(r *Resolver) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if initialized, exists := scope[v.name.Lexeme]; exists && !initialized {
			loxError(v.name, "Can't read local variable in its own initializer!")
		}
	}

	r.resolveLocal(v, v.name)
}
func (b Break) Resolve(r *Resolver) {
	if r.currentLoop == NoLoop {
		loxError(Token{Type: BREAK, Lexeme: "break"}, "Can't use 'break' outside of a loop.")
	}
}
func (r *Resolver) resolveFunction(function Function, ftype FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ftype

	r.beginScope()
	for _, param := range function.params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStatements(function.body)
	r.endScope()

	r.currentFunction = enclosingFunction
}
func (r *Resolver) resolveStatements(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveStatement(statement Stmt) {
	if resolvable, ok := statement.(Resolvable); ok {
		resolvable.Resolve(r)
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	if len(r.scopes) > 0 {
		r.scopes = r.scopes[:len(r.scopes)-1]
	}
}

func (r *Resolver) declare(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		loxError(name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true

}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			// r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
