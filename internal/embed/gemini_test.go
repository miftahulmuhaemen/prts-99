package embed

import (
	"context"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestNewEmbeddingModel(t *testing.T) {
	// Mock environment variables
	os.Setenv("GEMINI_API_KEY", "test-api-key")
	os.Setenv("GEMINI_MODEL", "test-model")

	// Create a context
	ctx := context.Background()

	// Mock the genai.Client
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		t.Fatalf("Failed to create genai client: %v", err)
	}

	// Create a GeminiClient
	geminiClient := &GeminiClient{
		Client: client,
	}

	// Call NewEmbeddingModel
	embeddingModel := geminiClient.NewEmbeddingModel()

	// Assertions
	assert.NotNil(t, embeddingModel)
	assert.NotNil(t, embeddingModel.EmbedModel)
	assert.NotNil(t, embeddingModel.QueryModel)
	assert.Equal(t, genai.TaskTypeRetrievalDocument, embeddingModel.EmbedModel.TaskType)
	assert.Equal(t, genai.TaskTypeRetrievalQuery, embeddingModel.QueryModel.TaskType)
}
