package lexer

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type lexType string

const (
	Literal   lexType = "literal"
	Number    lexType = "number"
	Separator lexType = "separator"
	Operator  lexType = "operator"
)

var (
	operators  = []string{"+", "-", "*", "/", "^", "%", "=", "!"}
	separators = []string{",", "(", ")"}

	// ErrSameTypeFollow is thrown when two numbers follow each others
	ErrSameTypeFollow = errors.New("sequence of two with exclusively numbers")
)

type Lexer struct {
	Type  lexType
	Value string
}

// Lex returns the Lexer of the content
func Lex(content string) ([]*Lexer, error) {
	var lexr []*Lexer
	for _, w := range strings.Split(content, " ") {
		word, err := lexWord(w)
		if err != nil {
			return nil, err
		}
		lexr = append(lexr, word...)
		for j := 0; j < len(lexr)-1; j++ {
			if lexr[j].Type == Number && lexr[j].Type == lexr[j+1].Type {
				return nil, errors.Join(
					ErrSameTypeFollow,
					fmt.Errorf(
						"not possible to have %s %s",
						lexr[j].Value,
						lexr[j+1].Value,
					),
				)
			}
		}
	}
	return lexr, nil
}

// lexWord returns the Lexer of the word
func lexWord(w string) ([]*Lexer, error) {
	if isDigit(w) {
		if []rune(w)[0] == '-' {
			return []*Lexer{
				{Operator, "-"},
				{Number, w[1:]},
			}, nil
		} else if []rune(w)[0] == '+' {
			return []*Lexer{
				{Operator, "+"},
				{Number, w[1:]},
			}, nil
		}
		return []*Lexer{
			{Number, w},
		}, nil
	}
	var lexers []*Lexer
	sel := ""
	var tpe lexType

	fnUpdate := func(typ lexType) {
		if tpe == typ {
			return
		}
		if tpe != "" {
			lexers = append(lexers, &Lexer{tpe, sel})
		}
		sel = ""
		tpe = typ
	}

	fnUpdateUnique := func(typ lexType) {
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

// isDigit checks if the string contains a digit
func isDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isOperator checks if the rune is an operator
func isOperator(s rune) bool {
	return slices.Contains(operators, string(s))
}

// isSeparator checks if the rune is a separator
func isSeparator(s rune) bool {
	return slices.Contains(separators, string(s))
}

func (l *Lexer) String() string {
	return fmt.Sprintf("%s('%s')", l.Type, l.Value)
}
