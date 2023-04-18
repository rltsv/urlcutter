package rest

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *HandlerShortener) *gin.Engine {
	router := gin.Default()
	router.Use(gzipHandler())

	router.POST("/", handler.CreateShortLink)
	router.POST("/api/shorten", handler.CreateShortLinkViaJSON)
	router.GET("/:id", handler.GetLinkByID)
	router.GET("/api/user/urls", handler.GetLinksByUser)

	return router
}
