package gomath

// Parse an expression with given Options
func Parse(expression string, opt *Options) (string, error) {
	lexed, err := lex(expression)
	if err != nil {
		return "", err
	}
	p, err := astParse(lexed, astTypeCalculation)
	if err != nil {
		return "", err
	}
	return p.Body.Eval(opt)
}
