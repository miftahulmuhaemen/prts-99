package handler

import (
	"chat-ak-wikia/internal/db"
	"chat-ak-wikia/internal/embed"
	"chat-ak-wikia/internal/service"
)

type Handler struct {
	user           service.User
	embeddingModel *embed.EmbeddingModel
	qdrant         *db.QdrantClient
}

func NewHandler(us service.User, em *embed.EmbeddingModel, qd *db.QdrantClient) *Handler {
	return &Handler{
		user:           us,
		embeddingModel: em,
		qdrant:         qd,
	}
}
