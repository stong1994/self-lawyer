package main

import (
	"io"
	"log"
	"net/http"
	"self-lawyer/chat"
	"self-lawyer/repo"
	"self-lawyer/vector"
)

func main() {
	log.Println("connecting ollama")
	ollama := vector.NewOllama()
	log.Println("connecting milvus")
	milvus := repo.NewMilvus(ollama)

	log.Println("starting chating server")
	chat := chat.NewOllama(milvus)

	mux := http.NewServeMux()
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		answer, err := chat.Complete(r.Context(), string(bytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(answer))
	})
	log.Println("listen server on :8888")
	http.ListenAndServe(":8888", mux)
}
