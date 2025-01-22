package lexer

import (
	"errors"
	"slices"
	"strconv"
	"strings"
)

const (
	Literal   string = "Literal"
	Number    string = "Number"
	Seperator string = "Seperator"
	Operator  string = "Operator"
)

var (
	InvalidSeparatorErr = errors.New("invalid separator")
	operators           = []string{"+", "-", "*", "/", "^", "%", "="}
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
	if strings.Contains(w, ",") {
		if strings.Contains(w[:len(w)-1], ",") {
			return nil, errors.Join(InvalidSeparatorErr, errors.New(w+" contains an invalid comma"))
		}
		return []*Lexer{
			{Literal, w[:len(w)-1]},
			{Seperator, ","},
		}, nil
	}
	if isOperator(w) {
		return []*Lexer{
			{Operator, w},
		}, nil
	}
	return []*Lexer{
		{Literal, w},
	}, nil
}

func isDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isOperator(s string) bool {
	return slices.Contains(operators, s)
}

func (l *Lexer) String() string {
	return l.Type + "(" + l.Value + ")"
}
