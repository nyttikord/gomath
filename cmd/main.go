package main

import (
	"flag"
	"fmt"
	"github.com/nyttikord/gomath"
	"os"
	"strings"
)

var (
	precision = uint(6)
)

func init() {
	flag.UintVar(&precision, "p", precision, "precision level")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("Usage: %s <subcommand>\nUse '%s help' for more information.\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}
	subcommand := args[0]
	switch subcommand {
	case "help":
		if len(args) != 1 {
			fmt.Printf("'help' does not take any arguments.\n")
			os.Exit(1)
		}
		fmt.Printf(
			"Usage: %s [flags] <subcommand>\n\nSubcommands:\n"+
				"- help               -> print this help text\n"+
				"- eval <expression>  -> evaluate an expression.\n"+
				"- latex <expression> -> convert an expression to LaTeX code.\n\n"+
				"Flags:\n"+
				"- p uint -> define the precision of the decimal approximation\n",
			os.Args[0],
		)
	case "eval":
		if len(args[1:]) == 0 {
			fmt.Printf("Usage: '%s eval <expression>'.\n", os.Args[0])
			os.Exit(1)
		}
		expression := strings.Join(args[1:], " ")
		res, err := gomath.Parse(expression)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Printf("Exact:   %s\n", res)
		fmt.Printf("Decimal: %s", res.Approx(int(precision)))
		if res.IsExact(int(precision)) {
			fmt.Printf(" (exact)")
		} else {
			fmt.Printf(" (not exact)")
		}
		fmt.Println()
	case "latex":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: '%s latex <expression>'.\n", os.Args[0])
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
		fmt.Printf("Unknown subcommand: %s\nUse '%s help' for more information.\n", subcommand, os.Args[0])
	}
}
