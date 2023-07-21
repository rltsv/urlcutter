package repository

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
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

func (s *MemoryStorage) SaveLink(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	for _, val := range s.Links {
		if val.UserID == dto.UserID && val.OriginalURL == dto.OriginalURL {
			return val.UserID, val.ShortURL, ErrLinkAlreadyExist
		}
	}

	s.Links = append(s.Links, dto)
	log.Println(s.Links)
	return dto.UserID, dto.ShortURL, nil
}

func (s *MemoryStorage) GetLink(ctx context.Context, dto entity.Link) (longurl string, err error) {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	if ok := s.CheckUserInMemory(dto); !ok {
		return "", ErrUserIsNotFound
	}
	for _, val := range s.Links {
		if val.LinkID == dto.LinkID && val.UserID == dto.UserID {
			return val.OriginalURL, nil
		}
	}
	return "", ErrLinkNotFound
}

func (s *MemoryStorage) GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error) {
	links = make([]entity.SendLinkDTO, 0)
	if ok := s.CheckUserInMemory(dto); !ok {
		return nil, ErrUserIsNotFound
	} else {
		for _, val := range s.Links {
			if val.UserID == dto.UserID {
				link := entity.SendLinkDTO{
					ShortURL:    val.ShortURL,
					OriginalURL: val.OriginalURL,
				}
				links = append(links, link)
			}
		}
	}
	return links, nil
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

func (s *MemoryStorage) Ping(ctx context.Context) error {
	return errors.New("there is no management system for db in this configuration")
}
