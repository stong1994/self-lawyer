package repo

import (
	"context"
	"fmt"
	"log"
	"self-lawyer/document_parser"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Vector interface {
	Embed(ctx context.Context, content string) ([]float32, error)
	GetDim() int
}

type Milvus struct {
	client client.Client
	vector Vector
}

func NewMilvus(vector Vector) *Milvus {
	// Connect to Milvus
	client := GetClient(context.Background())
	m := &Milvus{
		client: client,
		vector: vector,
	}
	m.InitCollection(context.Background())
	return m
}

func (m *Milvus) Store(ctx context.Context, laws document_parser.Laws) error {
	var (
		titles     []string
		contents   []string
		embeddings [][]float32
	)
	for _, law := range laws {
		for _, content := range law.Content {
			titles = append(titles, law.Title)
			contents = append(contents, content)
			embedding, err := m.vector.Embed(ctx, content)
			if err != nil {
				return err
			}
			embeddings = append(embeddings, embedding)
		}
	}
	// Insert vector to Milvus
	rst, err := m.client.Insert(
		ctx,
		collectionName,
		"",
		entity.NewColumnVarChar(titleCol, titles),
		entity.NewColumnVarChar(contentCol, contents),
		entity.NewColumnFloatVector(embeddingCol, m.vector.GetDim(), embeddings),
	)
	if err != nil {
		return err
	}
	log.Printf("Inserted %d rows\n", rst.Len())
	return m.client.Flush(ctx, collectionName, false)
}

type ContentResult struct {
	Content string
	ID      int64
}

type SearchResult struct {
	Title   string
	Content []ContentResult
}

func (s SearchResult) GetContents() []string {
	cs := make([]string, 0, len(s.Content))
	for _, c := range s.Content {
		cs = append(cs, c.Content)
	}
	return cs
}

type SearchResults []SearchResult

func (s SearchResults) Print() {
	for _, result := range s {
		fmt.Println(result.Title)
		for _, content := range result.Content {
			fmt.Printf("\tid: %d, content: %s\n", content.ID, content.Content)
		}
	}
}

func (m *Milvus) Search(ctx context.Context, content string) (SearchResults, error) {
	// Embed content
	embedding, err := m.vector.Embed(ctx, content)
	if err != nil {
		return nil, err
	}
	sp, err := entity.NewIndexHNSWSearchParam(16)
	if err != nil {
		return nil, err
	}
	// Search similar vectors
	res, err := m.client.Search(
		ctx,
		collectionName,
		nil,
		"",
		[]string{idCol, titleCol, contentCol, embeddingCol},
		[]entity.Vector{entity.FloatVector(embedding)},
		embeddingCol,
		entity.L2,
		16,
		sp,
	)
	if err != nil {
		return nil, err
	}
	var searchResult []SearchResult
	for _, row := range res {
		id := row.Fields.GetColumn(idCol)
		title := row.Fields.GetColumn(titleCol)
		content := row.Fields.GetColumn(contentCol)
		for i := 0; i < title.Len(); i++ {
			d, err := id.GetAsInt64(i)
			if err != nil {
				return nil, err
			}
			t, err := title.GetAsString(i)
			if err != nil {
				return nil, err
			}
			c, err := content.GetAsString(i)
			if err != nil {
				return nil, err
			}
			if len(searchResult) > 0 && searchResult[len(searchResult)-1].Title == t {
				searchResult[len(searchResult)-1].Content = append(searchResult[len(searchResult)-1].Content, ContentResult{Content: c, ID: d})
			} else {
				searchResult = append(searchResult, SearchResult{
					Title:   t,
					Content: []ContentResult{{Content: c, ID: d}},
				})
			}
		}
	}
	return searchResult, nil
}
