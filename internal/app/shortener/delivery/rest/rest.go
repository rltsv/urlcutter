package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"io"
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	if len(respBody) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	shortLink := hs.useCase.CreateShortLink(ctx, string(respBody))
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.Write([]byte(shortLink))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (hs *HandlerShortener) CreateShortLinkViaJSON(c *gin.Context) {
	ctx := c.Request.Context()

	rawValue, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer c.Request.Body.Close()

	ValueIn := &entity.InputData{}
	var shortLink string

	if err := json.Unmarshal(rawValue, &ValueIn); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	shortLink = hs.useCase.CreateShortLink(ctx, ValueIn.URL)

	ValueOut := entity.OutputData{
		Response: shortLink,
	}

	rawShortLink, err := json.Marshal(ValueOut)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.Writer.Header().Set("content-type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)
	_, err = c.Writer.Write(rawShortLink)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (hs *HandlerShortener) GetLinkByID(c *gin.Context) {

	ctx := c.Request.Context()

	id, _ := strconv.Atoi(c.Param("id"))

	origLink, err := hs.useCase.GetLinkByID(ctx, id)
	if err == repository.ErrLinkNotFound {
		c.AbortWithStatus(http.StatusBadRequest)
		c.Error(err)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, origLink)
	}
}
