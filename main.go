package main

import (
	"context"
	"fmt"
	"self-lawyer/document_parser"
	"self-lawyer/repo"
	"self-lawyer/vector"
)

func main() {
	laws, err := document_parser.Parse()
	if err != nil {
		panic(err)
	}
	// laws.Print()
	_ = laws
	fmt.Println("got laws ", len(laws))
	ollama, err := vector.NewOllama()
	if err != nil {
		panic(err)
	}
	milvus := repo.NewMilvus(ollama)

	findContent, err := milvus.Search(context.Background(), "最低工资")
	if err != nil {
		panic(err)
	}
	findContent.Print()
}
