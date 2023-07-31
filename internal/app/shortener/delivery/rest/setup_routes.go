package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/middleware"
)

func SetupRouter(handler *HandlerShortener) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.Gzip(), middleware.CheckCookie())

	router.POST("/", handler.CreateShortLink)
	router.POST("/api/shorten", handler.CreateShortLinkViaJSON)
	router.GET("/:id", handler.GetLinkByID)
	router.GET("/api/user/urls", handler.GetLinksByUser)

	return router
}
