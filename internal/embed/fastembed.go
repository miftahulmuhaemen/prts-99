package embed

import (
	"github.com/anush008/fastembed-go"
)

func GeneratePassageEmbeddings(documents []string) [][]float32 {
	// With custom options
	options := fastembed.InitOptions{
		Model:     fastembed.BGEBaseENV15,
		CacheDir:  "model_cache",
		MaxLength: 512,
	}

	model, err := fastembed.NewFlagEmbedding(&options)
	if err != nil {
		panic(err)
	}
	defer model.Destroy()

	// Generate embeddings with a batch-size of 25, defaults to 256
	embeddings, err := model.PassageEmbed(documents, 25) //  -> Embeddings length: 4
	if err != nil {
		panic(err)
	}

	return embeddings // Use embeddings as needed
}

func GenerateQueryEmbedding(query string) ([]float32, error) {
	// With custom options
	options := fastembed.InitOptions{
		Model:     fastembed.BGEBaseENV15,
		CacheDir:  "model_cache",
		MaxLength: 512,
	}

	model, err := fastembed.NewFlagEmbedding(&options)
	if err != nil {
		panic(err)
	}
	defer model.Destroy()

	// Generate embeddings with a batch-size of 25, defaults to 256
	embeddings, err := model.QueryEmbed(query)

	return embeddings, err
}
