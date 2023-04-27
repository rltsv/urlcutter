package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

func CheckCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userid string

		cookie, err := c.Request.Cookie("token")
		if err != nil && err == http.ErrNoCookie {
			userid = hex.EncodeToString(GenerateUserID())
		} else if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			userid = auth.DecryptToken(cookie)
		}

		switch c.ContentType() {
		case "text/plain":
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "failed while read body",
				})
			}
			dto := entity.CreateLinkDTO{
				UserID:  userid,
				LongURL: string(body),
			}

			dtoBytes, err := json.Marshal(&dto)
			if err != nil {
				log.Print(err)
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(dtoBytes))

		case "application/json":
			var dto entity.CreateLinkDTO

			if err = c.ShouldBindJSON(&dto); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "failed while parse json",
				})
			}

			dto.UserID = userid
			dtoBytes, err := json.Marshal(&dto)
			if err != nil {
				log.Print(err)
			}
			c.Request.Body = io.NopCloser(bytes.NewReader(dtoBytes))
		}

		c.Next()

	}
}

// GenerateUserID generate userID
func GenerateUserID() []byte {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("error while generateUserID: %v\n", err)
		return nil
	}
	return b
}
