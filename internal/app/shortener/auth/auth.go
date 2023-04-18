package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"net/http"
	"net/url"
)

var SecretKey = []byte("e24dac5623a329bdf879b21f56866830")

// CreateToken create token to set in cookie
func CreateToken(userID string) (token []byte) {
	src := []byte(userID)

	aesblock, err := aes.NewCipher(SecretKey)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	nonce := SecretKey[len(SecretKey)-aesgcm.NonceSize():]

	token = aesgcm.Seal(nil, nonce, src, nil)

	return token
}

// DecryptToken get out encrypted user id
func DecryptToken(token *http.Cookie) (userID string) {
	tokenValue := token.Value
	value, err := url.QueryUnescape(tokenValue)

	aesblock, err := aes.NewCipher(SecretKey)
	if err != nil {
		fmt.Printf("error2: %v\n", err)
		return
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error3: %v\n", err)
		return
	}

	nonce := SecretKey[len(SecretKey)-aesgcm.NonceSize():]

	userid, err := aesgcm.Open(nil, nonce, []byte(value), nil)
	if err != nil {
		fmt.Printf("error4: %v\n", err)
		return
	}
	return string(userid)
}
