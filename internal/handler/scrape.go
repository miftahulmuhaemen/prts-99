package handler

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"

	"github.com/qdrant/go-client/qdrant"

	internal "chat-ak-wikia/internal/scrapper"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Scrape(c echo.Context) error {
	operators, err := internal.Scrapper(5, "https://arknights.wiki.gg/wiki/Operator/6-star")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error scraping data")
	}

	points := []*qdrant.PointStruct{}
	for _, operator := range operators {
		values, payloads := operator.Values()
		for valueIndex, value := range values {
			// Embedding operator value
			embedded := h.embeddingModel.Embed(context.TODO(), strings.Join(value[1:], " "))
			vector := qdrant.NewVectors(embedded...)

			// Hashing operator name and value to get unique id
			hash := fnv.New64a()
			hash.Write(fmt.Appendf(nil, "%s%s", operator.OperatorName, value[0]))
			id := hash.Sum64()

			points = append(points, &qdrant.PointStruct{
				Id:      qdrant.NewIDNum(id),
				Vectors: vector,
				Payload: payloads[valueIndex],
			})
		}
	}

	_, err = h.qdrant.Update(context.TODO(), points)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, operators)
}
