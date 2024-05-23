package repo

import (
	"context"
	"log"
	"self-lawyer/document_parser"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

const (
	// Milvus instance proxy address, may verify in your env/settings
	milvusAddr = `localhost:19530`

	dbName                                    = "self_lawyer"
	collectionName                            = `laws`
	idCol, titleCol, contentCol, embeddingCol = "id", "title", "content", "embedding"
)

func GetClient(ctx context.Context) client.Client {
	// setup context for client creation, use 10 seconds here
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{
		Address: milvusAddr,
	})
	if err != nil {
		// handling error and exit, to make example simple here
		log.Fatal("failed to connect to milvus:", err.Error())
	}
	return c
}

// basic milvus operation example
func (m *Milvus) InitCollection(ctx context.Context) {
	dbs, err := m.client.ListDatabases(ctx)
	if err != nil {
		log.Fatal("failed to list databases:", err.Error())
	}
	exists := false
	for _, db := range dbs {
		if db.Name == "self_lawyer" {
			exists = true
		}
	}
	if !exists {
		err = m.client.CreateDatabase(ctx, "self_lawyer")
		if err != nil {
			log.Fatal("failed to create database:", err.Error())
		}
	}
	err = m.client.UsingDatabase(ctx, "self_lawyer")
	if err != nil {
		log.Fatal("failed to use database:", err.Error())
	}

	// first, lets check the collection exists
	collExists, err := m.client.HasCollection(ctx, collectionName)
	if err != nil {
		log.Fatal("failed to check collection exists:", err.Error())
	}
	if collExists {
		if err = m.client.LoadCollection(ctx, collectionName, false); err != nil {
			log.Fatal("failed to load collection:", err.Error())
		}
		return
	}
	log.Println("initializing collection...")

	// define collection schema
	schema := entity.NewSchema().WithName(collectionName).WithDescription("law data collection").
		// currently primary key field is compulsory, and only int64 is allowed
		WithField(entity.NewField().WithName(idCol).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)).
		// title and content
		WithField(entity.NewField().WithName(titleCol).WithDataType(entity.FieldTypeVarChar).WithMaxLength(50)).
		WithField(entity.NewField().WithName(contentCol).WithDataType(entity.FieldTypeVarChar).WithMaxLength(1024)).
		// also the vector field is needed
		WithField(entity.NewField().WithName(embeddingCol).WithDataType(entity.FieldTypeFloatVector).WithDim(int64(m.vector.GetDim())))

	err = m.client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}
	params := map[string]string{
		"M":              "16",
		"efConstruction": "100",
		"ef":             "20",
		"metric_type":    "L2",
	}
	err = m.client.CreateIndex(ctx, collectionName, embeddingCol, entity.NewGenericIndex("idx_embedding", entity.HNSW, params), false)
	if err != nil {
		log.Fatal("failed to create index:", err.Error())
	}
	if err = m.client.LoadCollection(ctx, collectionName, false); err != nil {
		log.Fatal("failed to load collection:", err.Error())
	}
	laws, err := document_parser.Parse()
	if err != nil {
		log.Fatal("failed to parse laws:", err.Error())
	}
	if err = m.Store(ctx, laws); err != nil {
		log.Fatal("failed to fill data:", err.Error())
	}
}
