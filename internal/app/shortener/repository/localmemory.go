package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"log"
	"sync"
)

var (
	ErrLinkAlreadyExist = errors.New("this link already shortened")
	ErrUnknownLink      = errors.New("unknown link")
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
	linkID := string(GenerateLinkID())
	if dto.UserID != "" {
		return &entity.Link{
			LinkID:   linkID,
			UserID:   dto.UserID,
			LongURL:  dto.LongURL,
			ShortURL: fmt.Sprintf("%s/%s", baseurl, linkID),
		}
	} else {
		return &entity.Link{
			LinkID:   linkID,
			UserID:   string(GenerateUserID()),
			LongURL:  dto.LongURL,
			ShortURL: fmt.Sprintf("%s/%s", baseurl, linkID),
		}
	}

}

func (s *MemoryStorage) SaveLinkInMemoryStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	if dto.UserID == "" {
		link := CreateNewLink(s.AppConfig.BaseURL, dto)
		s.Links = append(s.Links, *link)
		return link.UserID, link.ShortURL, nil
	} else { // dto.UserID != ""
		for _, val := range s.Links {
			if val.LongURL == dto.LongURL && val.UserID == dto.UserID {
				return "", "", ErrLinkAlreadyExist
			} else {
				link := CreateNewLink(s.AppConfig.BaseURL, dto)
				s.Links = append(s.Links, *link)
				return link.UserID, link.ShortURL, nil
			}
		}
	}
	return "", "", nil
}

func (s *MemoryStorage) GetLinkFromInMemoryStorage(ctx context.Context, dto entity.Link) (longurl string, err error) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	for _, val := range s.Links {
		if val.LinkID == dto.LinkID && val.UserID == dto.UserID {
			return val.LongURL, nil
		} else {
			return "", ErrUnknownLink
		}
	}
	return "", nil
}

func (s *MemoryStorage) CheckLinkInMemoryStorage(ctx context.Context, linkdto entity.Link) (id int, err error) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return 0, nil
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
