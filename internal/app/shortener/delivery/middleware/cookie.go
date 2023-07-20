package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rltsv/urlcutter/internal/app/shortener/auth"
)

type CookieContextKey string

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
		key := CookieContextKey("userid")
		ctx := context.WithValue(c.Request.Context(), key, userid)
		ctxWithVal := c.Request.WithContext(ctx)
		c.Request = ctxWithVal

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
