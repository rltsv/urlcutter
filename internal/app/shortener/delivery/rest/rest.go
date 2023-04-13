package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/rltsv/urlcutter/internal/app/shortener/usecase/shortener"
	"io"
	"net/http"
	"strconv"
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
	err := json.NewDecoder(c.Request.Body).Decode(&dto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
			"error":   err.Error(),
		})
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
			"error":   err.Error(),
		})
	}

	cookie, err := c.Request.Cookie("token")
	if err != nil {
		switch err {
		case http.ErrNoCookie:
			userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "invalid data",
					"error":   err.Error(),
				})
			}
			token := auth.CreateToken(userid)

			c.Writer.WriteHeader(http.StatusCreated)
			c.SetCookie("token", string(token), 0, "", "", false, false)
			_, err = c.Writer.Write([]byte(shorturl))
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		userid := auth.DecryptToken(cookie)
		dto.UserID = userid

		userid, shorturl, err := hs.useCase.CreateShortLink(ctx, dto)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "invalid data",
				"error":   err.Error(),
			})
		}

		token := auth.CreateToken(userid)

		c.Writer.WriteHeader(http.StatusCreated)
		c.SetCookie("token", string(token), 0, "", "", false, false)
		_, err = c.Writer.Write([]byte(shorturl))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
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

	id, _ := strconv.Atoi(c.Param("id"))

	origLink, err := hs.useCase.GetLinkByID(ctx, id)
	if err != nil && err == repository.ErrLinkNotFound {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, origLink)
	}
}
