package main

import (
	"fmt"
	"self-lawyer/document_parser"
)

func main() {
	laws, err := document_parser.Parse()
	if err != nil {
		panic(err)
	}
	_ = laws
	fmt.Println("==================================================")
	laws.Print()
}
