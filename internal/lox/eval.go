package lox

type Eval interface {
	Eval() (interface{}, error)
}

func (l *Literal) Eval() (interface{}, error) {
	return l.value, nil
}

func (b *Binary) Eval() (interface{}, error) {
	return nil, nil
}
