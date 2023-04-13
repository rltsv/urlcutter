package auth

import (
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreateToken(t *testing.T) {
	type request struct {
		userID    string
		secretKey []byte
	}

	tests := []struct {
		name    string
		request request
	}{
		{
			name: "check token generate",
			request: request{
				userID:    string(repository.GenerateUserID()),
				secretKey: SecretKey,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := CreateToken(tc.request.userID)

			cookie := &http.Cookie{Value: string(token)}

			userID := DecryptToken(cookie)
			assert.Equal(t, tc.request.userID, userID)
		})
	}
}
