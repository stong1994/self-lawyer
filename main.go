package main

import (
	"context"
	"fmt"
	"log"
	"self-lawyer/chat"
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
	ollama := vector.NewOllama()
	milvus := repo.NewMilvus(ollama)

	chat := chat.NewOllama(milvus)
	answer, err := chat.Complete(context.Background(), "公司发放的工资低于最低工资怎么办？")
	if err != nil {
		panic(err)
	}
	log.Println(answer)
}
