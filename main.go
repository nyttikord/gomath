package gomath

// ParseAndCalculate an expression with given Options
func ParseAndCalculate(expression string, opt *Options) (string, error) {
	return parseAndEval(expression, opt, astTypeCalculation)
}

// ParseAndConvertToLatex an expression with given Options
func ParseAndConvertToLatex(expression string, opt *Options) (string, error) {
	return parseAndEval(expression, opt, astTypeLatex)
}

func parseAndEval(expression string, opt *Options, tpe astType) (string, error) {
	lexed, err := lex(expression)
	if err != nil {
		return "", err
	}
	p, err := astParse(lexed, tpe)
	if err != nil {
		return "", err
	}
	return p.Body.Eval(opt)
}
