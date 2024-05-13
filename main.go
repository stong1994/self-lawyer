package main

import (
	"context"
	"fmt"
	"self-lawyer/document_parser"
	"self-lawyer/search_engine"
)

func main() {
	laws, err := document_parser.Parse()
	if err != nil {
		panic(err)
	}
	laws.Print()
	_ = laws
	searchEngine, err := search_engine.NewOllama()
	if err != nil {
		panic(err)
	}
	vector, err := searchEngine.Embed(context.Background(), "hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(vector)
}
