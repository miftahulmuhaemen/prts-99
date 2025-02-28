// Create basic web server with echo framework
// and return Hello World! as response
package main

import (
	"chat-ak-wikia/internal/db"
	em "chat-ak-wikia/internal/embed_model"
	internal "chat-ak-wikia/internal/scrapper"
	"os"
	"os/signal"
	"strings"
	"time"

	"context"
	"fmt"
	"hash/fnv"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/qdrant/go-client/qdrant"
)

func main() {

	// Not override the existing environment variables
	_ = godotenv.Load()

	// Create new qdrant instance
	qClient := db.NewClient(context.Background())
	defer qClient.Close()

	// Create new Gemini instance
	gClient := em.NewClient(context.Background())
	model := gClient.NewEmbeddingModel()
	defer gClient.Close()

	// Create new echo instance
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	defer e.Close()

	e.GET("/query", func(c echo.Context) error {
		query := c.QueryParam("query")
		if query == "" {
			return c.String(http.StatusBadRequest, "Query is required")
		}

		embed := model.Query(context.TODO(), query)
		res, err := qClient.SearchPoints(context.TODO(), embed)
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
	})

	e.GET("/scrape", func(c echo.Context) error {
		collector := colly.NewCollector(
			colly.AllowedDomains("arknights.wiki.gg"),
			colly.CacheDir("./cache"),
		)

		operators, err := internal.Scrapper(5, "https://arknights.wiki.gg/wiki/Operator/6-star", collector)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error scraping data")
		}

		h := fnv.New64a()
		points := []*qdrant.PointStruct{}
		for _, operator := range operators {
			strs, ngl := operator.OperatorToStrings()
			print("Length: ", len(strs), "-", len(ngl), "\n")
			for strI, str := range strs {
				embedded := model.Embedding(context.TODO(), strings.Join(str[1:], " "))
				vector := qdrant.NewVectors(embedded...)

				h.Write(fmt.Appendf(nil, "%s%s", operator.OperatorName, str[0]))
				id := h.Sum64()

				print("Id: \n", id)

				points = append(points, &qdrant.PointStruct{
					Id:      qdrant.NewIDNum(id),
					Vectors: vector,
					Payload: ngl[strI],
				})
			}
		}

		_, err = qClient.Update(points)
		if err != nil {
			fmt.Println(err)
		}

		return c.JSON(http.StatusOK, operators)
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
