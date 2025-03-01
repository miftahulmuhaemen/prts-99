package handler

import (
	"chat-ak-wikia/internal/router/middleware"
	"chat-ak-wikia/internal/utils"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	guestUsers := v1.Group("/users")
	guestUsers.POST("", h.SignUp)
	guestUsers.POST("/login", h.Login)

	user := v1.Group("/user", jwtMiddleware)
	user.GET("", h.CurrentUser)
	user.PUT("", h.UpdateUser)

	profiles := v1.Group("/profiles", jwtMiddleware)
	profiles.GET("/:username", h.GetProfile)

	chat := v1.Group("/chat")
	chat.GET("", h.Query)

	scrape := v1.Group("/scrape")
	scrape.GET("", h.Scrape)
}
