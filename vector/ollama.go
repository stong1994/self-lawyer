package vector

import (
	"context"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	llm   *ollama.LLM
	model string
	dim   int
}

func (o *Ollama) GetDim() int {
	return o.dim
}

type Option func(*Ollama)

func WithOptionSetModel(model string) Option {
	return func(o *Ollama) {
		o.model = model
	}
}

func NewOllama(opts ...Option) *Ollama {
	o := new(Ollama)
	for _, opt := range opts {
		opt(o)
	}

	o.fillConfig()

	llm, err := ollama.New(ollama.WithModel(o.model))
	if err != nil {
		log.Fatal(err)
	}
	o.llm = llm

	o.dim = o.getEmbeddingDim()
	log.Printf("got embedding engine, model is %s, dim is %d", o.model, o.dim)
	return o
}

func (o *Ollama) fillConfig() {
	if o.model == "" {
		o.model = "llama3"
	}
}

func (o *Ollama) getEmbeddingDim() int {
	data, err := o.Embed(context.Background(), "hello")
	if err != nil {
		log.Fatal(err)
	}
	return len(data)
}

func (o *Ollama) Embed(ctx context.Context, content string) ([]float32, error) {
	result, err := o.llm.CreateEmbedding(ctx, []string{content})
	if err != nil {
		return nil, err
	}

	return result[0], err
}
