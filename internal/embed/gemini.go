package embed

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	*genai.Client
}

type EmbeddingModel struct {
	EmbedModel *genai.EmbeddingModel
	QueryModel *genai.EmbeddingModel
}

func NewClient(ctx context.Context) *GeminiClient {
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	return &GeminiClient{
		Client: client,
	}
}

func (c *GeminiClient) NewEmbeddingModel() *EmbeddingModel {
	emRet := c.EmbeddingModel(os.Getenv("GEMINI_MODEL"))
	emRet.TaskType = genai.TaskTypeRetrievalDocument
	emQry := c.EmbeddingModel(os.Getenv("GEMINI_MODEL"))
	emQry.TaskType = genai.TaskTypeRetrievalQuery
	return &EmbeddingModel{
		EmbedModel: emRet,
		QueryModel: emQry,
	}
}

func (c *EmbeddingModel) Embed(ctx context.Context, text string) []float32 {
	res, err := c.EmbedModel.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		log.Fatal(err)
	}
	return res.Embedding.Values
}

func (c *EmbeddingModel) Query(ctx context.Context, text string) []float32 {
	res, err := c.QueryModel.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		log.Fatal(err)
	}
	return res.Embedding.Values
}
