package lox

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func (lc LoxClass) toString() string {
	return lc.name
}

func (lc LoxClass) call(arguments []interface{}) interface{} {
	inst := &LoxInstance{class: &lc, fields: make(map[string]interface{})}
	initializer := lc.findMethod("init")
	if initializer != nil {
		initializer.bind(inst).call(arguments)
	}

	return inst
}

func (lc LoxClass) arity() int {
	initializer := lc.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

func (lc *LoxClass) findMethod(name string) *LoxFunction {
	if f, exists := lc.methods[name]; exists {
		return f
	}
	return nil
}
