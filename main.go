// Create basic web server with echo framework
// and return Hello World! as response
package main

import (
	"chat-ak-wikia/internal/db"
	"chat-ak-wikia/internal/embed_model"
	internal "chat-ak-wikia/internal/scrapper"
	"log"
	"strings"

	"context"
	"fmt"
	"hash/fnv"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/qdrant/go-client/qdrant"
)

func main() {

	// Not override the existing environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create new echo instance
	e := echo.New()

	// Create new qdrant instance
	vectorCollection := "arknights"
	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		panic(err)
	}

	// Check if collection already exist or not
	_, err = qdrantClient.GetCollectionInfo(context.Background(), vectorCollection)
	if err != nil {
		onDisk := true
		err = qdrantClient.CreateCollection(
			context.Background(),
			&qdrant.CreateCollection{
				CollectionName: vectorCollection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Cosine,
				}),
				OnDiskPayload: &onDisk,
			},
		)
		if err != nil {
			panic(err)
		}
	}

	e.GET("/query", func(c echo.Context) error {
		query := "Which faction Archetto is?"
		embed := embed_model.GenerateQuery(query)

		db.SearchPoints(context.Background(), qdrantClient, vectorCollection, embed)

		return c.JSON(http.StatusOK, []string{})
	})

	// Define route
	// e.GET("/", func(c echo.Context) error {
	// 	collector := colly.NewCollector(
	// 		colly.AllowedDomains("arknights.wiki.gg"),
	// 		colly.CacheDir("./cache"),
	// 	)

	// 	operators, err := internal.Scrapper(5, "https://arknights.wiki.gg/wiki/Operator/6-star", collector)
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, "Error scraping data")
	// 	}

	// 	// embeddings := [][][]float32{}
	// 	// for _, operator := range operators {
	// 	// 	embedded := embed_model.GeneratePassageEmbeddings(operator.OperatorToStrings())
	// 	// 	embeddings = append(embeddings, embedded);
	// 	// }

	// 	points := []*qdrant.PointStruct{}
	// 	for _, operator := range operators {
	// 		embedded := embed_model.GeneratePassageEmbeddings(operator.OperatorToStrings())

	// 		fmt.Println("Operator :", operator.OperatorName)
	// 		for _, embed := range embedded {
	// 			vector := qdrant.NewVectors(embed...)
	// 			n, _ := rand.Int(rand.Reader, big.NewInt(1000))
	// 			id := n.Int64() + 1

	// 			fmt.Println(id)

	// 			points = append(points, &qdrant.PointStruct{
	// 				Id: qdrant.NewIDNum(uint64(id)),
	// 				// TODO: Multivector
	// 				Vectors: vector,
	// 				Payload: qdrant.NewValueMap(map[string]any{
	// 					"operator_name": operator.OperatorName,
	// 				}),
	// 				// TODO: Payload: qdrant.NewValueMap(model.StructToMap(operator)),
	// 			})
	// 		}
	// 	}

	// 	onWait := true
	// 	update, err := qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
	// 		CollectionName: vectorCollection,
	// 		Points:         points,
	// 		Wait:           &onWait,
	// 	})

	// 	print(update)

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	return c.JSON(http.StatusOK, operators)
	// })

	e.GET("/scrape-gemini", func(c echo.Context) error {
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
				embedded := embed_model.GenerateEmbedding(strings.Join(str[1:], " "))
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

		onWait := true
		_, err = qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
			CollectionName: vectorCollection,
			Points:         points,
			Wait:           &onWait,
		})

		if err != nil {
			fmt.Println(err)
		}

		return c.JSON(http.StatusOK, operators)
	})

	// Start server
	e.Start(":8080")
}
