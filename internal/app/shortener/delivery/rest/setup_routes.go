package rest

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *HandlerShortener) *gin.Engine {
	router := gin.Default()

	router.POST("/", handler.CreateShortLink)
	router.GET("/:id", handler.GetLinkByID)

	return router
}
