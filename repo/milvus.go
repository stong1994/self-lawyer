package repo

import (
	"context"
	"self-lawyer/document_parser"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Vector interface {
	Embed(ctx context.Context, content string) ([]float32, error)
}

type Milvus struct {
	client client.Client
	vector Vector
}

func NewMilvus(vector Vector) *Milvus {
	// Connect to Milvus
	client := InitCollection(context.Background())
	return &Milvus{
		client: client,
		vector: vector,
	}
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
	_, err := m.client.Insert(
		ctx,
		collectionName,
		"",
		entity.NewColumnVarChar(titleCol, titles),
		entity.NewColumnVarChar(contentCol, contents),
		entity.NewColumnFloatVector(embeddingCol, dim, embeddings),
	)
	if err != nil {
		return err
	}
	return m.client.Flush(ctx, collectionName, false)
}

func (m *Milvus) Search(ctx context.Context, content string) (document_parser.Laws, error) {
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
		10,
		sp,
	)
	if err != nil {
		return nil, err
	}
	var laws []document_parser.Law
	for _, row := range res {
		title := row.Fields.GetColumn(titleCol)
		content := row.Fields.GetColumn(contentCol)
		for i := 0; i < title.Len(); i++ {
			t, err := title.GetAsString(i)
			if err != nil {
				return nil, err
			}
			c, err := content.GetAsString(i)
			if err != nil {
				return nil, err
			}
			laws = append(laws, document_parser.Law{
				Title:   t,
				Content: []string{c},
			})
		}
	}
	return laws, nil
}
