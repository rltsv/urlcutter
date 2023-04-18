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
		c.AbortWithError(http.StatusBadRequest, errors.New("failed while read body"))
		return
	}

	dto.LongURL = string(longURL)

	cookie, err := c.Request.Cookie("token")
	switch err {
	case http.ErrNoCookie:
		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("failed while create shorturl"))
			return
		}
		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie("token", string(auth.CreateToken(userid)), 0, "", "", false, false)
		c.Writer.Write([]byte(shorturl))
	default:
		dto.UserID = auth.DecryptToken(cookie)

		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil && err == repository.ErrLinkAlreadyExist {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie("token", string(auth.CreateToken(userid)), 0, "", "", false, false)
		c.Writer.Write([]byte(shorturl))
	}
}

func (hs *HandlerShortener) CreateShortLinkViaJSON(c *gin.Context) {
	//ctx := c.Request.Context()

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

	//shortLink = hs.useCase.CreateShortLink(ctx, ValueIn.URL)

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

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "there is no created links by this user",
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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
		log.Print(errors.New("error: problem with link id"))
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
			"message": "there is no user ",
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "there is no created links by this user",
			})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	dto := entity.GetAllLinksDTO{UserID: auth.DecryptToken(cookie)}
	links, err := hs.useCase.GetLinksByUser(ctx, dto)
	log.Print(links)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	body, err := json.Marshal(&links)

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Header().Set("content-type", "application/json")
	c.SetCookie("token", cookie.Value, 0, "", "", false, false)
	c.Writer.Write(body)

}
