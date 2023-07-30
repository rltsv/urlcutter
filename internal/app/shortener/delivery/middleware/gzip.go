package middleware

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomRW struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (c CustomRW) Write(data []byte) (int, error) {
	return c.buf.Write(data)
}

func NewCustomRW(ctx *gin.Context) *CustomRW {
	return &CustomRW{
		ResponseWriter: ctx.Writer,
		buf:            &bytes.Buffer{},
	}
}

func Gzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		NewCRW := NewCustomRW(c)
		c.Writer = NewCRW

		if c.Request.Header.Get(`Content-Encoding`) == `gzip` {
			gzReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "error occurred while create new gzip reader",
					"error":   err.Error(),
				})
				return
			}
			c.Request.Body = gzReader

			defer gzReader.Close()
		}

		c.Next()

		if strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
			NewCRW.Header().Add("Content-Encoding", "gzip")
			gzWriter := gzip.NewWriter(NewCRW.ResponseWriter)
			defer gzWriter.Close()
			gzWriter.Write(NewCRW.buf.Bytes())
		} else {
			NewCRW.ResponseWriter.Write(NewCRW.buf.Bytes())
		}
	}
}
