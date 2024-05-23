package vector

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	llm *ollama.LLM
	dim int
}

func (o *Ollama) GetDim() int {
	return o.dim
}

func NewOllama() *Ollama {
	// llm, err := ollama.New(ollama.WithModel("nomic-embed-text:v1.5"))
	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal(err)
	}

	o := &Ollama{
		llm: llm,
	}
	data, err := o.Embed(context.Background(), "hello")
	if err != nil {
		log.Fatal(err)
	}
	o.dim = len(data)
	return o
}

func (o *Ollama) Embed(ctx context.Context, content string) ([]float32, error) {
	result, err := o.llm.CreateEmbedding(ctx, []string{content})
	if err != nil {
		return nil, err
	}

	return result[0], err
}
