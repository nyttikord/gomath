package lexer

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

const (
	Literal   string = "literal"
	Number    string = "number"
	Separator string = "separator"
	Operator  string = "operator"
)

var (
	operators  = []string{"+", "-", "*", "/", "^", "%", "="}
	separators = []string{",", "(", ")"}
)

type Lexer struct {
	Type, Value string
}

func Lex(content []string) ([][]*Lexer, error) {
	var lexers [][]*Lexer
	for i, line := range content {
		if strings.HasPrefix(line, "#") {
			continue
		}
		var lexer []*Lexer
		for _, w := range strings.Split(line, " ") {
			word, err := lexWord(w)
			if err != nil {
				return nil, errors.Join(errors.New("line "+strconv.Itoa(i+1)+" has an error"), err)
			}
			lexer = append(lexer, word...)
		}
		lexers = append(lexers, lexer)
	}
	return lexers, nil
}

func lexWord(w string) ([]*Lexer, error) {
	if isDigit(w) {
		return []*Lexer{
			{Number, w},
		}, nil
	}
	var lexers []*Lexer
	sel := ""
	tpe := ""

	fnUpdate := func(typ string) {
		if tpe == typ {
			return
		}
		if tpe != "" {
			lexers = append(lexers, &Lexer{tpe, sel})
		}
		sel = ""
		tpe = typ
	}

	for _, c := range []rune(w) {
		if isDigit(string(c)) || (c == '.' && sel != "") {
			fnUpdate(Number)
		} else if isOperator(c) {
			fnUpdate(Operator)
		} else if isSeparator(c) {
			fnUpdate(Separator)
		}
		sel += string(c)
	}
	lexers = append(lexers, &Lexer{tpe, sel})
	return lexers, nil
}

func isDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s rune) bool {
	return slices.Contains(operators, string(s))
}

func isSeparator(s rune) bool {
	return slices.Contains(separators, string(s))
}

func (l *Lexer) String() string {
	return l.Type + "(" + l.Value + ")"
}
