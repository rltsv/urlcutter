package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"io"
	"log"
	"net/http"
	"strconv"
)

type HandlerShortener struct {
	useCase shortener.Usecase
}

func NewHandlerShortener(shortenerUseCase shortener.Usecase) *HandlerShortener {
	return &HandlerShortener{
		useCase: shortenerUseCase,
	}
}

func (hs *HandlerShortener) CreateShortLink(c *gin.Context) {

	ctx := c.Request.Context()

	respBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.Error(c.Writer, "specify the request", 400)
		return
	}
	defer c.Request.Body.Close()

	if len(respBody) == 0 {
		http.Error(c.Writer, "where is nothing to short, check body", 400)
		return
	}

	shortLink := hs.useCase.CreateShortLink(ctx, string(respBody))

	c.Writer.WriteHeader(201)
	_, err = c.Writer.Write([]byte(shortLink))
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
		log.Fatal("", err)
		return
	}
}

func (hs *HandlerShortener) GetLinkByID(c *gin.Context) {

	ctx := c.Request.Context()

	id, _ := strconv.Atoi(c.Param("id"))

	origLink, err := hs.useCase.GetLinkByID(ctx, id)
	if err == repository.ErrLinkNotFound {
		http.Error(c.Writer, "there is no any link by this id", 400)
	}

	c.Writer.Header().Set("Location", origLink)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
