package repository

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
)

var (
	ErrLinkNotFound   = errors.New("link not found")
	ErrStorageIsEmpty = errors.New("storage is empty")
)

type ShortenerRepo interface {
	SaveLink(ctx context.Context, longLink string) (IDCount int)
	GetLink(ctx context.Context, id int) (longLink string, err error)
}

type LinksRepository struct {
	Storage map[int]string
	IDCount int
	Mux     *sync.Mutex
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		Storage: make(map[int]string),
		IDCount: 0,
		Mux:     &sync.Mutex{},
	}
}
func (l *LinksRepository) SaveLink(ctx context.Context, longLink string) (IDCount int) {

	for key, val := range l.Storage {
		if strings.EqualFold(val, longLink) {
			return key
		}
	}

	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.IDCount++
	l.Storage[l.IDCount] = longLink

	return l.IDCount

}

func (l *LinksRepository) GetLink(ctx context.Context, id int) (longLink string, err error) {

	if len(l.Storage) == 0 {
		return "", ErrStorageIsEmpty
	}

	l.Mux.Lock()
	defer l.Mux.Unlock()
	if val, ok := l.Storage[id]; !ok {
		log.Print("ya tut")
		return "", ErrLinkNotFound
	} else {
		return val, nil
	}

}
