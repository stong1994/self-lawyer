package search_engine

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	llm *ollama.LLM
}

func NewOllama() (*Ollama, error) {
	llm, err := ollama.New(ollama.WithModel("nomic-embed-text:v1.5"))
	if err != nil {
		log.Fatal(err)
	}
	return &Ollama{
		llm: llm,
	}, nil
}

func (o *Ollama) Embed(ctx context.Context, content string) ([]float32, error) {
	result, err := o.llm.CreateEmbedding(ctx, []string{content})
	if err != nil {
		return nil, err
	}

	return result[0], err
}
