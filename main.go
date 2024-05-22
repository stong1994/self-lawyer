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
	ollama, err := vector.NewOllama()
	if err != nil {
		panic(err)
	}
	milvus := repo.NewMilvus(ollama)

	err = milvus.Store(context.Background(), laws)
	if err != nil {
		panic(err)
	}
	findContent, err := milvus.Search(context.Background(), "解除劳动合同")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", findContent)
}
