package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"log"
	"sync"
)

type MemoryStorage struct {
	Links     []entity.Link
	Mux       *sync.RWMutex
	AppConfig config.Config
}

func NewMemoryStorage(cfg config.Config) *MemoryStorage {
	return &MemoryStorage{
		Links:     make([]entity.Link, 0),
		Mux:       new(sync.RWMutex),
		AppConfig: cfg,
	}
}

// CreateNewLink create new instance of link
func CreateNewLink(baseurl string, dto entity.Link) *entity.Link {

	encodedLinkID := hex.EncodeToString(GenerateLinkID())
	switch {
	case dto.UserID != "":
		return &entity.Link{
			LinkID:   encodedLinkID,
			UserID:   dto.UserID,
			LongURL:  dto.LongURL,
			ShortURL: fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
		}
	default:
		encodedUserID := hex.EncodeToString(GenerateUserID())
		return &entity.Link{
			LinkID:   encodedLinkID,
			UserID:   encodedUserID,
			LongURL:  dto.LongURL,
			ShortURL: fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
		}
	}
}

func (s *MemoryStorage) SaveLinkInMemoryStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	switch s.CheckUserInMemory(dto) {
	case false:
		link := CreateNewLink(s.AppConfig.BaseURL, dto)
		s.Links = append(s.Links, *link)
		return link.UserID, link.ShortURL, nil
	case true:
		for _, val := range s.Links {
			if val.UserID == dto.UserID && val.LinkID == dto.LinkID {
				return val.UserID, val.ShortURL, ErrLinkAlreadyExist
			} else {
				link := CreateNewLink(s.AppConfig.BaseURL, dto)
				s.Links = append(s.Links, *link)
				return link.UserID, link.ShortURL, nil
			}
		}
	}
	return userid, shorturl, err
}

func (s *MemoryStorage) GetLinkFromInMemoryStorage(ctx context.Context, dto entity.Link) (longurl string, err error) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	if ok := s.CheckUserInMemory(dto); !ok {
		return "", ErrUserIsNotFound
	}
	for _, val := range s.Links {
		if val.LinkID == dto.LinkID && val.UserID == dto.UserID {
			return val.LongURL, nil
		}
	}
	return "", ErrLinkNotFound
}

// CheckUserInMemory check are user in already in memory or not
func (s *MemoryStorage) CheckUserInMemory(dto entity.Link) (ok bool) {
	if dto.UserID == "" {
		return false
	}
	for _, val := range s.Links {
		if val.UserID == dto.UserID {
			return true
		}
	}
	return false
}

// GenerateLinkID generate LinkID
func GenerateLinkID() []byte {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("error while generateLinkID: %v\n", err)
		return nil
	}
	return b
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
