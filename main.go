package main

import (
	"context"
	"log"
	"self-lawyer/chat"
	"self-lawyer/repo"
	"self-lawyer/vector"
)

func main() {
	ollama := vector.NewOllama()
	milvus := repo.NewMilvus(ollama)

	chat := chat.NewOllama(milvus)
	answer, err := chat.Complete(context.Background(), "公司发放的工资低于最低工资怎么办？")
	if err != nil {
		panic(err)
	}
	log.Println(answer)
}
