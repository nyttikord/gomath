package main

import (
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
}
