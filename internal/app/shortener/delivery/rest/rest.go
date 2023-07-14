package rest

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
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

	dto.UserID = c.Request.Context().Value("userid").(string)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while read request body",
			"error":   err.Error(),
		})
		return
	}
	// format raw data from request body to string type and assign it to dto field
	dto.OriginalURL = string(rawBody)

	userid, shortLink, err := hs.useCase.CreateShortLink(ctx, dto)
	if err != nil && err == repository.ErrLinkAlreadyExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "this link already shortened",
		})
		return
	}
	bytesOut := struct {
		ShortLink string `json:"shortlink"`
	}{
		ShortLink: shortLink,
	}

	shortLinkBytes, err := json.Marshal(bytesOut)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while marshal string",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
	c.Writer.Write(shortLinkBytes)
}

func (hs *HandlerShortener) CreateShortLinkViaJSON(c *gin.Context) {
	ctx := c.Request.Context()
	var dto entity.CreateLinkDTO

	dto.UserID = c.Request.Context().Value("userid").(string)

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userid, shortlink, err := hs.useCase.CreateShortLink(ctx, dto)
	if err != nil && err == repository.ErrLinkAlreadyExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "this link already shortened",
		})
		return
	}
	bytesOut := struct {
		ShortLink string `json:"shortlink"`
	}{
		ShortLink: shortlink,
	}

	shortLinkBytes, err := json.Marshal(bytesOut)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while marshal string",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
	c.Writer.Write(shortLinkBytes)
}

func (hs *HandlerShortener) GetLinkByID(c *gin.Context) {
	ctx := c.Request.Context()

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while grab cookie from request",
		})
		return
	}

	userID := auth.DecryptToken(cookie)
	linkID := c.Param("id")
	if linkID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "check id path of url",
		})
	}

	dto := entity.GetLinkDTO{
		UserID: userID,
		LinkID: linkID,
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

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "there is no shortened links by this user",
			})
			return
		default:
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
	userid := auth.DecryptToken(cookie)

	dto := entity.GetAllLinksDTO{UserID: userid}
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
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
	c.Writer.Write(linksBytes)
}

func (hs *HandlerShortener) Ping(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	err := hs.useCase.Ping(ctx)
	if err != nil && err.Error() == "there is no management system for db in this configuration" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while pinging database",
			"error":   err.Error(),
		})
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed while pinging database",
			"error":   err.Error(),
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func (hs *HandlerShortener) BatchShortener(c *gin.Context) {
	ctx := c.Request.Context()

	var request []entity.CreateLinkDTO

	err := c.ShouldBind(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed while parse request body",
			"error":   err.Error(),
		})
		return
	}

	userid := c.Request.Context().Value("userid").(string)

	for idx, _ := range request {
		request[idx].UserID = userid
	}

	listOfShortUrls, err := hs.useCase.BatchShortener(ctx, request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
	c.JSON(http.StatusOK, listOfShortUrls)
}
