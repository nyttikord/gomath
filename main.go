package gomath

func Parse(content string, opt *Options) {
	lexed, err := lex(content)
	if err != nil {
		panic(err)
	}
	//for _, l := range lexed {
	//	s := ""
	//	for _, v := range l {
	//		s += v.String() + " "
	//	}
	//	println(s[:len(s)-1])
	//}
	p, err := astParse(lexed)
	if err != nil {
		panic(err)
	}
	//m, err := json.MarshalIndent(p, "", "  ")
	//if err != nil {
	//	panic(err)
	//}
	//println(string(m))
	for _, stmt := range p.Body {
		err = stmt.Eval(opt)
		if err != nil {
			panic(err)
		}
	}
}
