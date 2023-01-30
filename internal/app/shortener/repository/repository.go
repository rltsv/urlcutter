package repository

import (
	"context"
	"errors"
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
	Mutex   *sync.Mutex
	RwMutex *sync.RWMutex
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		Storage: make(map[int]string),
		IDCount: 0,
		Mutex:   new(sync.Mutex),
		RwMutex: new(sync.RWMutex),
	}
}
func (l *LinksRepository) SaveLink(ctx context.Context, longLink string) (IDCount int) {

	l.RwMutex.Lock()
	defer l.RwMutex.Unlock()
	for key, val := range l.Storage {
		if strings.EqualFold(val, longLink) {
			return key
		}
	}

	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.IDCount++
	l.Storage[l.IDCount] = longLink
	return l.IDCount

}

func (l *LinksRepository) GetLink(ctx context.Context, id int) (longLink string, err error) {
	l.RwMutex.RLock()
	defer l.RwMutex.RUnlock()

	if len(l.Storage) == 0 {
		return "", ErrStorageIsEmpty
	}

	if val, ok := l.Storage[id]; !ok {
		return "", ErrLinkNotFound
	} else {
		return val, nil

	}
}
