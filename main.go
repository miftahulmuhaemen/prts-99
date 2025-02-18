// Create basic web server with echo framework
// and return Hello World! as response
package main

import (
	internal "chat-ak-wikia/internal/scrapper"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/labstack/echo/v4"
)

func main() {
	// Create new echo instance
	e := echo.New()

	// Define route
	e.GET("/", func(c echo.Context) error {
		collector := colly.NewCollector(
			colly.AllowedDomains("arknights.wiki.gg"),
			colly.CacheDir("./cache"),
		)
		operator, err := internal.Scrapper(1, "https://arknights.wiki.gg/wiki/Operator/6-star", collector)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error scraping data")
		}
		return c.JSON(http.StatusOK, operator)
	})

	// Start server
	e.Start(":8080")
}
