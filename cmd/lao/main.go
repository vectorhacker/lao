package main

import (
	"io"
	"log"
	"os"

	"github.com/vectorhacker/lao/pkg/lao"
)

func main() {
	var r io.Reader
	{
		if len(os.Args) != 2 {
			r = os.Stdin
		} else {

			filePath := os.Args[1]
			f, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			r = f
		}
	}

	var tokenizer lao.Tokenizer
	{
		tokenizer = lao.NewTokenizer(r)
	}

	var parser lao.Parser
	{
		parser = lao.NewParser(tokenizer)
	}

	var interpreter lao.Interpreter
	{
		interpreter = lao.NewInterpreter(os.Stdout)
	}

	statements, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	err = interpreter.Execute(statements)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}
}
