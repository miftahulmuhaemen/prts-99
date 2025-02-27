package embed_model

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateEmbedding(text string) []float32 {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	em := client.EmbeddingModel("text-embedding-004")
	em.TaskType = genai.TaskTypeRetrievalDocument
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		log.Fatal(err)
	}
	return res.Embedding.Values
}

func GenerateQuery(text string) []float32 {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	em := client.EmbeddingModel("text-embedding-004")
	em.TaskType = genai.TaskTypeRetrievalQuery
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		log.Fatal(err)
	}
	return res.Embedding.Values
}
