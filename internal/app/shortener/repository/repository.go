package repository

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrLinkNotFound = errors.New("link not found")
)

type ShortenerRepo interface {
	CreateLink(ctx context.Context, longLink string) (id int)
	GetLinkByID(ctx context.Context, id int) (longLink string, err error)
	CheckLinkInMemory(longLink string) (id int)
}

type LinksRepository struct {
	Storage map[int]string
	IDCount int
	Mux     *sync.RWMutex
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		Storage: make(map[int]string),
		IDCount: 0,
		Mux:     new(sync.RWMutex),
	}
}
func (l *LinksRepository) CreateLink(ctx context.Context, longLink string) (id int) {

	id = l.CheckLinkInMemory(longLink)
	if id != 0 {
		return id
	}

	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.IDCount++
	l.Storage[l.IDCount] = longLink

	return l.IDCount
}

func (l *LinksRepository) GetLinkByID(ctx context.Context, id int) (longLink string, err error) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()
	if val, ok := l.Storage[id]; !ok {
		return "", ErrLinkNotFound
	} else {
		return val, nil
	}
}

func (l *LinksRepository) CheckLinkInMemory(longLink string) (id int) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()
	for key, val := range l.Storage {
		if strings.EqualFold(val, longLink) {
			return key
		}
	}
	return 0
}
