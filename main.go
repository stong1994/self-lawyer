package main

import (
	"encoding/json"
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
	serve(chat)
}

type chatRequest struct {
	Question string `json:"question"`
}

type chatResponse struct {
	Answer string `json:"answer"`
}

func serve(chat *chat.Ollama) {
	mux := http.NewServeMux()
	mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		var req chatRequest
		if err = json.Unmarshal(bytes, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		answer, err := chat.Complete(r.Context(), req.Question)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		res, _ := json.Marshal(chatResponse{Answer: answer})
		w.Write(res)
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
