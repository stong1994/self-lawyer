package chat

import (
	"context"
	"fmt"
	"log"
	"self-lawyer/repo"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var systemMessages = llms.TextParts(llms.ChatMessageTypeSystem, `
You are an AI lawyer assistant.
'system' will provide multi law items, you should filter the most related laws to the user's question. 
Your output must be in Chinese. If you are uncertain or the answer is not 
explicitly written in the documentation, please respond with "I'm sorry, I cannot assist with this.`)

type SearchEngine interface {
	Search(ctx context.Context, content string) (repo.SearchResults, error)
}

type Option func(*Ollama)

func WithModel(model string) Option {
	return func(o *Ollama) {
		o.model = model
	}
}

type Ollama struct {
	model        string
	llm          *ollama.LLM
	searchEngine SearchEngine
	// chatHistory  []llms.MessageContent
}

func NewOllama(searchEngine SearchEngine, opts ...Option) *Ollama {
	o := &Ollama{
		searchEngine: searchEngine,
		// chatHistory:  []llms.MessageContent{systemMessages},
	}
	for _, opt := range opts {
		opt(o)
	}
	o.fillDefaultConfig()

	llm, err := ollama.New(ollama.WithModel(o.model))
	if err != nil {
		log.Fatal(err)
	}
	o.llm = llm
	log.Printf("got chat engine, model is %s", o.model)

	return o
}

func (o *Ollama) fillDefaultConfig() {
	if o.model == "" {
		o.model = "llama3"
	}
}

func msgSendToAI(sr repo.SearchResults) string {
	var contents []string
	for _, chapter := range sr {
		for _, item := range chapter.Content {
			contents = append(contents, fmt.Sprintf("%s %s", chapter.Chapter, item.Content))
		}
	}
	return strings.Join(contents, " ")
}

func (o *Ollama) jointUserMessage(problem string, relatedLaws repo.SearchResults) []llms.MessageContent {
	return []llms.MessageContent{
		systemMessages,
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextPart(msgSendToAI(relatedLaws))},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextPart(problem)},
		},
	}
}

func (o *Ollama) Complete(ctx context.Context, problem string, writer func(chunk []byte) error) error {
	laws, err := o.searchEngine.Search(ctx, problem)
	if err != nil {
		return err
	}
	messages := o.jointUserMessage(problem, laws)
	log.Printf("completing, request message: %+v", messages)
	log.Print("completing, please wait...")
	res, err := o.llm.GenerateContent(
		ctx,
		messages,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			err = writer(chunk)
			if err != nil {
				return err
			}
			log.Println("streaming chunk: ", string(chunk))
			return nil
		}))
	if err != nil {
		return err
	}

	log.Print("got choices")
	for i, choice := range res.Choices {
		log.Printf("\tchoice %d: content: %s, stop reason: %s", i, choice.Content, choice.StopReason)
	}
	// TODO: return error
	return nil
}
