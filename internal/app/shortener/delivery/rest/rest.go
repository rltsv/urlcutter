package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
	"github.com/rltsv/urlcutter/internal/app/shortener/delivery/middleware"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
)

type HandlerShortener struct {
	useCase shortener.UsecaseShortener
}

func NewHandlerShortener(shortenerUseCase shortener.UsecaseShortener) *HandlerShortener {
	return &HandlerShortener{
		useCase: shortenerUseCase,
	}
}

func (hs *HandlerShortener) CreateShortLink(c *gin.Context) {
	ctx := c.Request.Context()
	var dto entity.CreateLinkDTO

	dto.UserID = c.Request.Context().Value(middleware.ContextKey).(string)

	if rawBody, err := io.ReadAll(c.Request.Body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while read request body",
		})
		return
	} else {
		dto.LongURL = string(rawBody)
	}

	userid, shortLink, err := hs.useCase.CreateShortLink(ctx, dto)
	if err != nil && err == repository.ErrLinkAlreadyExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "link already shortened",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "text/html")
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "/", "", false, false)
	c.Writer.Write([]byte(shortLink))
}

func (hs *HandlerShortener) CreateShortLinkViaJSON(c *gin.Context) {
	ctx := c.Request.Context()
	var dto entity.CreateLinkDTO

	dto.UserID = c.Request.Context().Value(middleware.ContextKey).(string)

	if rawBody, err := io.ReadAll(c.Request.Body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while read request body",
		})
		return
	} else {
		if err = json.Unmarshal(rawBody, &dto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "failed while unmarshal data",
			})
			log.Print(err)
			return
		}
	}

	userid, shortlink, err := hs.useCase.CreateShortLink(ctx, dto)
	if err != nil && err == repository.ErrLinkAlreadyExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "link already shortened",
		})
		log.Print(err)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "text/html")
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "/", "", false, false)
	c.Writer.Write([]byte(shortlink))
}

func (hs *HandlerShortener) GetLinkByID(c *gin.Context) {
	ctx := c.Request.Context()
	dto := entity.GetLinkDTO{}

	dto.UserID = c.Request.Context().Value(middleware.ContextKey).(string)

	linkID := c.Param("id")
	if linkID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "check id path of url",
		})
	} else {
		dto.LinkID = linkID
	}

	longLink, err := hs.useCase.GetLinkByUserID(ctx, dto)
	if err != nil && err == repository.ErrLinkNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "there is no any link by this id",
		})
		log.Print(err.Error())
	} else if err != nil && err == repository.ErrUserIsNotFound {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "there is no shortened links by this user",
		})
		log.Print(err.Error())
	} else {
		c.Redirect(http.StatusTemporaryRedirect, longLink)
	}
}

func (hs *HandlerShortener) GetLinksByUser(c *gin.Context) {
	ctx := c.Request.Context()
	dto := entity.GetAllLinksDTO{}

	dto.UserID = c.Request.Context().Value(middleware.ContextKey).(string)

	links, err := hs.useCase.GetLinksByUser(ctx, dto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "there is no shortened links by this user",
		})
		log.Print(err)
		return
	}

	linksBytes, err := json.Marshal(&links)

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	c.SetCookie("token", string(auth.CreateToken(dto.UserID)), 3600, "/", "", false, false)
	c.Writer.Write(linksBytes)
}
