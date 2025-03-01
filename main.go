// Create basic web server with echo framework
// and return Hello World! as response
package main

import (
	"chat-ak-wikia/internal/db"
	"chat-ak-wikia/internal/embed"
	"chat-ak-wikia/internal/handler"
	"chat-ak-wikia/internal/router"
	"chat-ak-wikia/internal/store"

	"context"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// Not override the existing environment variables
	_ = godotenv.Load()

	// Initiate Context
	ctx := context.Background()

	// Create new qdrant instance
	qClient := db.NewClient(
		ctx,
		os.Getenv("QDRANT_HOST"),
		os.Getenv("QDRANT_PORT"),
		os.Getenv("QDRANT_COLLECTION"),
	)
	qClient.AutoMigrate(ctx, true)

	// Create new DB instance
	d := db.New()
	db.AutoMigrate(d)

	// Create new Gemini instance
	gClient := embed.NewClient(ctx)
	model := gClient.NewEmbeddingModel()

	// Create new echo instance
	e := router.New()

	v1 := e.Group("/api/v1")

	us := store.NewUserStore(d)
	h := handler.NewHandler(us, model, qClient)
	h.Register(v1)

	e.Logger.Fatal(e.Start(":8080"))
}
