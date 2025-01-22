package main

import (
	"encoding/json"
	"flag"
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
	lexed, err := lex(content)
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
}
