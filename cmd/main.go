package main

import (
	"fmt"
	"github.com/nyttikord/gomath"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <subcommand>\nUse %s help for more information.\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}
	subcommand := os.Args[1]
	switch subcommand {
	case "help":
		if len(os.Args) != 2 {
			fmt.Printf("Help does not take any arguments.\n")
			os.Exit(1)
		}
		fmt.Printf(
			"Usage: %s <subcommand>\nSubcommands\n"+
				"- help -> print this help text\n"+
				"- eval <expression> -> evaluate an expression.\n"+
				"- latex <expression> -> convert an expression to LaTeX code.\n",
			os.Args[0],
		)
	case "eval":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s eval <expression>.\n", os.Args[0])
			os.Exit(1)
		}
		expression := strings.Join(os.Args[2:], " ")
		res, err := gomath.Parse(expression)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(res)
	case "latex":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: %s latex <expression>.\n", os.Args[0])
			os.Exit(1)
		}
		expression := strings.Join(os.Args[2:], " ")
		res, err := gomath.ParseAndConvertToLaTeX(expression, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(res)
	default:
		fmt.Printf("Unknown subcommand: %s\nUse %s help to more information.\n", subcommand, os.Args[0])
	}
}
