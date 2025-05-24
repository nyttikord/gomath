package gomath

import (
	"errors"
	"fmt"
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
	operators  = []string{"+", "-", "*", "/", "^", "%", "=", "{", "}"}
	separators = []string{",", "(", ")"}

	SameTypeFollowErr = errors.New("sequence of two with exclusively numbers")
)

type Lexer struct {
	Type, Value string
}

func Lex(content string) ([]*Lexer, error) {
	var lexer []*Lexer
	for _, w := range strings.Split(content, " ") {
		word, err := lexWord(w)
		if err != nil {
			return nil, err
		}
		lexer = append(lexer, word...)
		for j := 0; j < len(lexer)-1; j++ {
			if lexer[j].Type == Number && lexer[j].Type == lexer[j+1].Type {
				return nil, errors.Join(
					SameTypeFollowErr,
					fmt.Errorf(
						"not possible to have %s %s",
						lexer[j].Value,
						lexer[j+1].Value,
					),
				)
			}
		}
	}
	return lexer, nil
}

func lexWord(w string) ([]*Lexer, error) {
	if isDigit(w) {
		if []rune(w)[0] == '-' {
			return []*Lexer{
				{Operator, "-"},
				{Number, w[1:]},
			}, nil
		} else if []rune(w)[0] == '+' {
			return []*Lexer{
				{Number, w[1:]},
			}, nil
		}
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

	fnUpdateUnique := func(typ string) {
		if tpe == typ && sel == "" {
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
			fnUpdateUnique(Operator)
		} else if isSeparator(c) {
			fnUpdateUnique(Separator)
		} else {
			fnUpdate(Literal)
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
