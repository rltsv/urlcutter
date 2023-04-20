package rest

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"io"
	"log"
	"net/http"
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

	longURL, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.New("failed while read body"))
		return
	}
	defer c.Request.Body.Close()

	dto.LongURL = string(longURL)

	cookie, err := c.Request.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil && err == repository.ErrLinkAlreadyExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "link already shortened",
			})
			return
		}
		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
		c.Writer.Write([]byte(shorturl))
	default:
		dto.UserID = auth.DecryptToken(cookie)

		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil && err == repository.ErrLinkAlreadyExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "link already shortened",
			})
			return
		}

		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
		c.Writer.Write([]byte(shorturl))
	}
}

func (hs *HandlerShortener) CreateShortLinkViaJSON(c *gin.Context) {
	ctx := c.Request.Context()

	var dto entity.CreateLinkDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cookie, err := c.Request.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		userid, shorturl, _ := hs.useCase.CreateShortLink(ctx, dto)
		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie(
			"token",
			string(auth.CreateToken(userid)),
			3600,
			"",
			"",
			false,
			false,
		)
		c.Writer.Write([]byte(shorturl))
	default:
		dto.UserID = auth.DecryptToken(cookie)

		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil && err == repository.ErrLinkAlreadyExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "link already shortened",
			})
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		c.SetCookie(
			"token",
			string(auth.CreateToken(userid)),
			3600,
			"",
			"",
			false,
			false,
		)
		c.Writer.WriteHeader(http.StatusCreated)
		c.Writer.Write([]byte(shorturl))
	}
}

func (hs *HandlerShortener) GetLinkByID(c *gin.Context) {
	ctx := c.Request.Context()

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "there is no created links by this user",
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
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

	longURL, err := hs.useCase.GetLinkByUserID(ctx, dto)
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
		c.Redirect(http.StatusTemporaryRedirect, longURL)
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
		return
	}

	body, err := json.Marshal(&links)

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	c.SetCookie("token", string(auth.CreateToken(userid)), 3600, "", "", false, false)
	c.Writer.Write(body)
}
