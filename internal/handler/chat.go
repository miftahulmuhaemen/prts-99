package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Query(c echo.Context) error {
	query := c.QueryParam("query")
	if query == "" {
		return c.String(http.StatusBadRequest, "Query is required")
	}

	embed := h.embeddingModel.Query(context.TODO(), query)
	res, err := h.qdrant.SearchPoints(context.TODO(), embed)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	ret := []string{}
	for _, r := range res {
		fmt.Println(r.Id, r.Score)
		if payload := r.GetPayload(); payload != nil {
			ret = append(ret, fmt.Sprintf("%v", payload))
		}
	}

	return c.JSON(http.StatusOK, ret)
}
