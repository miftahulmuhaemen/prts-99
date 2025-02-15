// Create basic web server with echo framework 
// and return Hello World! as response
package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"chat-ak-wikia/pkg/internal"
)

func main() {
	// Create new echo instance
	e := echo.New()

	// Define route
	e.GET("/", func(c echo.Context) error {
		// call PrintOperator function from internal package
		operator, err := scrapper.Scrapper()
        if err != nil {
            return c.String(http.StatusInternalServerError, "Error scraping data")
        }
        scrapper.PrintOperator(operator)
		return c.String(http.StatusOK, "Hello World!")
	})

	// Start server
	e.Start(":8080")
}