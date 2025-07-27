package lexer

type TokenList struct {
	index int
	list  []*Lexer
}

func (list *TokenList) Current() *Lexer {
	if list.index < 0 || list.Empty() {
		return nil
	}
	return list.list[list.index]
}

func (list *TokenList) Next() bool {
	list.index++
	return !list.Empty()
}

func (list *TokenList) Empty() bool {
	return list.index >= len(list.list)
}

func (list *TokenList) String() string {
	s := "["
	for _, l := range list.list {
		s += l.String() + " "
	}
	return s[:len(s)-1] + "]"
}
