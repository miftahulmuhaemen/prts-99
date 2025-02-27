package db

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

func SearchPoints(ctx context.Context, client *qdrant.Client, collectionName string, queryVector []float32) {
	// Define search parameters
	limit := uint64(5) // Number of results to return
	// threshold := float32(0.5) // Similarity threshold

	// Create the search request
	searchRequest := &qdrant.QueryPoints{
		CollectionName: collectionName,
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
	client.Query(ctx, searchRequest)
	response, err := client.Query(ctx, searchRequest)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Process the search results
	for _, result := range response {
		log.Printf("Point ID: %v, Score: %f", result.GetId(), result.GetScore())
		// Access payload if WithPayload is true
		if payload := result.GetPayload(); payload != nil {
			log.Printf("Payload: %v", payload)
		}
		// Access vector if WithVectors is true
		if vector := result.GetVectors(); vector != nil {
			log.Printf("Vector: %v", vector)
		}
	}

}
