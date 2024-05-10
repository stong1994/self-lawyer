package main

import "self-laywer/document_parser"

func main() {
	laws, err := document_parser.Parse()
	if err != nil {
		panic(err)
	}
	laws.Print()
	_ = laws
}
