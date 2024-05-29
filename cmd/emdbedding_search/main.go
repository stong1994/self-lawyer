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
	question             = ""
	reset                = false
	embedding_cache      = false
	embedding_cache_path = ""
)

func main() {
	flag.BoolVar(&reset, "reset", false, "reset the database")
	flag.StringVar(&question, "question", "试用期最长几个月?", "reset the database")
	flag.BoolVar(&embedding_cache, "embedding_cache", false, "get/set embedding with cache")
	flag.StringVar(&embedding_cache_path, "embedding_cache_path", "embedding_cache.text", "specify the path of embedding_cache")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("connecting ollama")
	var ollama repo.Vector
	ollama = vector.NewOllama(vector.WithOptionSetModel(os.Getenv("EMBEDDING_MODEL")))
	if embedding_cache {
		ollama = vector.NewCacheOllama(
			ollama.(*vector.Ollama),
			vector.NewCache(vector.WithCacheOptionSetCachePath(embedding_cache_path)))
	}
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
