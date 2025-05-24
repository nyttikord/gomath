package gomath

func Parse(content string, opt *Options) (string, error) {
	lexed, err := lex(content)
	if err != nil {
		panic(err)
	}
	p, err := astParse(lexed, "return")
	if err != nil {
		panic(err)
	}
	return p.Body.Eval(opt)
}
