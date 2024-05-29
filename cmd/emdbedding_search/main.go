package main

import (
	"context"
	"flag"
	"log"
	"os"
	"self-lawyer/repo"
	"self-lawyer/vector"

	"github.com/joho/godotenv"
)

var (
	question = ""
	reset    = false
)

func main() {
	flag.BoolVar(&reset, "reset", false, "reset the database")
	flag.StringVar(&question, "question", "试用期最长几个月?", "reset the database")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("connecting ollama")
	ollama := vector.NewOllama(vector.OptionSetModel(os.Getenv("EMBEDDING_MODEL")))
	log.Println("connecting milvus")
	milvus := repo.NewMilvus(ollama)

	if reset {
		log.Println("reseting system")
		milvus.DropDatabase(context.Background())
		milvus.InitCollection(context.Background())
	}

	log.Println("searching question: ", question)
	rst, err := milvus.Search(context.Background(), question)
	if err != nil {
		panic(err)
	}
	log.Println("search result: ")
	rst.Print()
}
