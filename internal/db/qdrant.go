package db

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type Client struct {
	*qdrant.Client
}

var vectorCollection = os.Getenv("QDRANT_COLLECTION")

func NewClient(ctx context.Context) *Client {

	port, err := strconv.Atoi(os.Getenv("QDRANT_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_HOST"),
		Port: port,
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = client.GetCollectionInfo(context.Background(), vectorCollection)
	if err != nil {
		onDisk := true
		err = client.CreateCollection(
			context.Background(),
			&qdrant.CreateCollection{
				CollectionName: vectorCollection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Cosine,
				}),
				OnDiskPayload: &onDisk,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	return &Client{
		Client: client,
	}
}

func (c *Client) Update(points []*qdrant.PointStruct) (*qdrant.UpdateResult, error) {
	onWait := true
	return c.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: vectorCollection,
		Points:         points,
		Wait:           &onWait,
	})
}

func (c *Client) SearchPoints(ctx context.Context, queryVector []float32) ([]*qdrant.ScoredPoint, error) {
	// Define search parameters
	limit := uint64(5) // Number of results to return
	// threshold := float32(0.5) // Similarity threshold

	// Create the search request
	searchRequest := &qdrant.QueryPoints{
		CollectionName: vectorCollection,
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          &limit,
		// WithPayload:    qdrant.NewWithPayload(true),
		// WithVectors:    qdrant.NewWithVectors(true),
		// ScoreThreshold: &threshold,
		// Optional parameters:
		// Filter:        &qdrant.Filter{}, // To filter results based on conditions
		// Params:        &qdrant.SearchParams{}, // Additional search parameters
		// WithPayload:   true, // Include payload in the response
		// WithVectors:   true, // Include vectors in the response
		// ScoreThreshold: 0.5, // Minimum score threshold for results
	}

	// Execute the search
	c.Query(ctx, searchRequest)
	response, err := c.Query(ctx, searchRequest)

	return response, err

	// Process the search results
	// for _, result := range response {
	// 	log.Printf("Point ID: %v, Score: %f", result.GetId(), result.GetScore())
	// 	// Access payload if WithPayload is true
	// 	if payload := result.GetPayload(); payload != nil {
	// 		log.Printf("Payload: %v", payload)
	// 	}
	// 	// Access vector if WithVectors is true
	// 	if vector := result.GetVectors(); vector != nil {
	// 		log.Printf("Vector: %v", vector)
	// 	}
	// }
}
