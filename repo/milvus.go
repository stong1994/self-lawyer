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

func (m *Milvus) Store(ctx context.Context, laws []document_parser.Laws) error {
	var (
		kinds      []string
		chapters   []string
		items      []string
		embeddings [][]float32
	)
	for _, law := range laws {
		for _, chapter := range law.Chapters {
			for _, item := range chapter.Items {
				kinds = append(kinds, law.Kind)
				chapters = append(chapters, chapter.Chapter)
				items = append(items, item.Content)
				embedding, err := m.vector.Embed(ctx, item.Content)
				if err != nil {
					return err
				}
				embeddings = append(embeddings, embedding)
			}
		}
	}
	// Insert vector to Milvus
	rst, err := m.client.Insert(
		ctx,
		collectionName,
		"",
		entity.NewColumnVarChar(kindCol, kinds),
		entity.NewColumnVarChar(chapterCol, chapters),
		entity.NewColumnVarChar(contentCol, items),
		entity.NewColumnFloatVector(embeddingCol, m.vector.GetDim(), embeddings),
	)
	if err != nil {
		return err
	}
	log.Printf("Inserted %d rows\n", rst.Len())
	return m.client.Flush(ctx, collectionName, false)
}

type ContentResult struct {
	Content  string
	ID       int64
	Distance float32
}

type SearchResult struct {
	Kind    string
	Chapter string
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
		fmt.Printf("%s %s", result.Kind, result.Chapter)
		for _, content := range result.Content {
			fmt.Printf("\tid: %d, socre: %f, content: %s\n", content.ID, content.Distance, content.Content)
		}
	}
}

func (m *Milvus) Search(ctx context.Context, content string) (SearchResults, error) {
	// Embed content
	embedding, err := m.vector.Embed(ctx, content)
	if err != nil {
		return nil, err
	}
	sp, err := entity.NewIndexHNSWSearchParam(10)
	if err != nil {
		return nil, err
	}
	// Search similar vectors
	res, err := m.client.Search(
		ctx,
		collectionName,
		nil,
		"",
		[]string{idCol, chapterCol, contentCol, embeddingCol},
		[]entity.Vector{entity.FloatVector(embedding)},
		embeddingCol,
		entity.COSINE,
		10,
		sp,
	)
	if err != nil {
		return nil, err
	}
	var searchResult []SearchResult
	for _, row := range res {
		id := row.Fields.GetColumn(idCol)
		chapter := row.Fields.GetColumn(chapterCol)
		kind := row.Fields.GetColumn(kindCol)
		content := row.Fields.GetColumn(contentCol)
		for i := 0; i < chapter.Len(); i++ {
			d, err := id.GetAsInt64(i)
			if err != nil {
				return nil, err
			}
			t, err := chapter.GetAsString(i)
			if err != nil {
				return nil, err
			}
			c, err := content.GetAsString(i)
			if err != nil {
				return nil, err
			}
			k, err := kind.GetAsString(i)
			if err != nil {
				return nil, err
			}
			if len(searchResult) > 0 && searchResult[len(searchResult)-1].Chapter == t {
				searchResult[len(searchResult)-1].Content = append(searchResult[len(searchResult)-1].Content, ContentResult{Content: c, ID: d, Distance: row.Scores[i]})
			} else {
				searchResult = append(searchResult, SearchResult{
					Kind:    k,
					Chapter: t,
					Content: []ContentResult{{Content: c, ID: d, Distance: row.Scores[i]}},
				})
			}
		}
	}
	return searchResult, nil
}
