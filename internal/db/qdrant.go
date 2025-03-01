package db

import (
	"context"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	*qdrant.Client
	collection string
}

func NewClient(ctx context.Context, host string, port string, collection string) *QdrantClient {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                   host,
		Port:                   portInt,
		SkipCompatibilityCheck: true,
	})
	if err != nil {
		panic(err)
	}
	return &QdrantClient{
		Client:     client,
		collection: collection,
	}
}

func (c *QdrantClient) AutoMigrate(ctx context.Context, onDisk bool) error {
	isExist, err := c.CollectionExists(ctx, c.collection)
	if err != nil {
		panic(err)
	}
	if !isExist {
		return c.CreateCollection(
			ctx,
			&qdrant.CreateCollection{
				CollectionName: c.collection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Cosine,
				}),
				OnDiskPayload: &onDisk,
			},
		)
	}
	return nil
}

func (c *QdrantClient) Update(ctx context.Context, points []*qdrant.PointStruct) (*qdrant.UpdateResult, error) {
	onWait := true
	return c.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: c.collection,
		Points:         points,
		Wait:           &onWait,
	})
}

func (c *QdrantClient) SearchPoints(ctx context.Context, queryVector []float32) ([]*qdrant.ScoredPoint, error) {
	// Define search parameters
	limit := uint64(5) // Number of results to return
	// threshold := float32(0.5) // Similarity threshold

	// Create the search request
	searchRequest := &qdrant.QueryPoints{
		CollectionName: c.collection,
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
		// WithVectors:    qdrant.NewWithVectors(true),
		// ScoreThreshold: &threshold,
		// Optional parameters:
		// Filter:        &qdrant.Filter{}, // To filter results based on conditions
		// Params:        &qdrant.SearchParams{}, // Additional search parameters
	}

	// Execute the search
	c.Query(ctx, searchRequest)
	response, err := c.Query(ctx, searchRequest)

	return response, err
}
