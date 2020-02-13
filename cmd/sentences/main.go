package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vectorhacker/lao/pkg/lao"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	p := lao.NewParser(lao.NewTokenizer(f))

	nodes, err := p.Parse()

	for _, node := range nodes {
		switch statement := node.(type) {
		case lao.RemStatement:

			fmt.Printf("REM Statement: %s\n", func() string {
				s := ""
				for _, t := range statement.Tokens() {
					s += t.Value + " "
				}
				return s
			}())
		case lao.PrintStatement:

			fmt.Printf(
				"PRINT Statement: %s %s\n",
				statement.Tokens()[0].Value,
				func() string {
					if len(statement.Tokens()) == 1 {
						return ""
					}
					return statement.Tokens()[1].Value
				}(),
			)
		case lao.AssignmentStatement:
			fmt.Printf("Asignment statement: %s\n", strings.Join(func() []string {
				s := []string{}

				for _, token := range statement.Tokens() {
					s = append(s, token.Value)
				}

				return s
			}(), " "))
		case lao.EndStatement:
			fmt.Printf("END Statement: %s%s\n", statement.Tokens()[0].Value, statement.Tokens()[1].Value)
		case lao.ReadStatement:
			fmt.Printf("READ Statment: %s %s\n", statement.Tokens()[0].Value, statement.Tokens()[1].Value)
		case lao.IfStatement:
			fmt.Printf("IF Satement: %s\n", strings.Join(func() []string {
				s := []string{}

				for _, token := range statement.Tokens() {
					s = append(s, token.Value)
				}

				return s
			}(), " "))
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
