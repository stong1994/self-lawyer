package chat

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"self-lawyer/repo"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var systemMessages = llms.TextParts(llms.ChatMessageTypeSystem, `
You are an AI lawyer assistant.
When asked for your name, you must respond with "Lawyer Copilot".
'system' will present a legal situation for which you will provide advice and relevant legal provisions. 
Please only provide advice related to this situation. Based on the specific sections from the documentation, 
answer the question only using that information. Please be aware that if there are any updates to the legal provisions, 
please reference the most current content. Your output must be in Chinese. If you are uncertain or the answer is not 
explicitly written in the documentation, please respond with "I'm sorry, I cannot assist with this.`)

type SearchEngine interface {
	Search(ctx context.Context, content string) (repo.SearchResults, error)
}

type Ollama struct {
	llm          *ollama.LLM
	searchEngine SearchEngine
	chatHistory  []llms.MessageContent
}

func NewOllama(searchEngine SearchEngine) *Ollama {
	// llm, err := ollama.New(ollama.WithModel("nomic-embed-text:v1.5"))
	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal(err)
	}

	o := &Ollama{
		llm:          llm,
		searchEngine: searchEngine,
		chatHistory:  []llms.MessageContent{systemMessages},
	}
	return o
}

func (o *Ollama) jointUserMessage(problem string, relatedLaws repo.SearchResults) []llms.MessageContent {
	o.chatHistory = append(o.chatHistory, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextPart(problem)},
	})
	text, _ := json.Marshal(relatedLaws)
	systemMessage := llms.MessageContent{
		Role:  llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextPart(string(text))},
	}
	o.chatHistory = append(o.chatHistory, systemMessage)

	return o.chatHistory
}

func (o *Ollama) jointAIMessage(answer string) []llms.MessageContent {
	o.chatHistory = append(o.chatHistory, llms.MessageContent{
		Role:  llms.ChatMessageTypeAI,
		Parts: []llms.ContentPart{llms.TextPart(answer)},
	})
	return o.chatHistory
}

func (o *Ollama) Complete(ctx context.Context, problem string, writer io.Writer) error {
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
			data := []byte("data: " + string(chunk) + "\n\n")
			_, err = writer.Write(data)
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
	if len(res.Choices) > 0 {
		o.jointAIMessage(res.Choices[0].Content)
		return nil
	}
	// TODO: return error
	return nil
}
