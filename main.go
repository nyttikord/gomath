package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/anhgelus/gomath/lexer"
	"os"
	"strings"
)

var (
	input string
)

func init() {
	flag.StringVar(&input, "i", "", "input file")
}

func main() {
	flag.Parse()
	b, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}
	content := strings.Split(string(b), "\n")
	lexed, err := lexer.Lex(content)
	if err != nil {
		panic(err)
	}
	for _, l := range lexed {
		s := ""
		for _, v := range l {
			s += v.String() + " "
		}
		println(s[:len(s)-1])
	}
	p, err := parse(lexed)
	if err != nil {
		panic(err)
	}
	m, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		panic(err)
	}
	println(string(m))
	for _, stmt := range p.Body {
		f, err := stmt.Eval()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s - %f\n", f.String(), f.Float())
	}
}
