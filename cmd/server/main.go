package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"self-lawyer/chat"
	"self-lawyer/repo"
	"self-lawyer/vector"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("connecting ollama")
	ollama := vector.NewOllama(vector.OptionSetModel(os.Getenv("EMBEDDING_MODEL")))
	log.Println("connecting milvus")
	milvus := repo.NewMilvus(ollama)

	log.Println("starting chating server")
	chat := chat.NewOllama(milvus, chat.OptionSetModel(os.Getenv("COMPLETING_MODEL")))
	serve(chat, milvus)
}

type chatRequest struct {
	Question string `json:"question"`
}

func serve(chat *chat.Ollama, milvus *repo.Milvus) {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		// Set the necessary headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		log.Println("question", string(bytes))
		var req chatRequest
		if err = json.Unmarshal(bytes, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = chat.Complete(r.Context(), req.Question, func(chunk []byte) error {
			fmt.Fprintf(w, "data: %s\n\n", string(chunk))
			flusher, ok := w.(http.Flusher)
			if !ok {
				// The ResponseWriter doesn't support the Flusher interface, so we can't stream
				return nil
			}
			flusher.Flush()
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
	mux.HandleFunc("/reset_all", func(w http.ResponseWriter, r *http.Request) {
		milvus.DropDatabase(r.Context())
		milvus.InitCollection(r.Context())
		w.Write([]byte("ok"))
	})

	log.Println("listen server on :8888")
	http.ListenAndServe(":8888", corsMiddleware(mux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
