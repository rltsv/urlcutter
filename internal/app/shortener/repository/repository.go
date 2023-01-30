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
	mutex   *sync.Mutex
	rwMutex *sync.RWMutex
}

func NewLinksRepository() *LinksRepository {
	return &LinksRepository{
		Storage: make(map[int]string),
		IDCount: 0,
		mutex:   new(sync.Mutex),
		rwMutex: new(sync.RWMutex),
	}
}
func (l *LinksRepository) SaveLink(ctx context.Context, longLink string) (IDCount int) {

	l.rwMutex.Lock()
	defer l.rwMutex.Unlock()
	for key, val := range l.Storage {
		if strings.EqualFold(val, longLink) {
			return key
		}
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.IDCount++
	l.Storage[l.IDCount] = longLink
	return l.IDCount

}

func (l *LinksRepository) GetLink(ctx context.Context, id int) (longLink string, err error) {
	l.rwMutex.RLock()
	defer l.rwMutex.RUnlock()

	if len(l.Storage) == 0 {
		return "", ErrStorageIsEmpty
	}

	if val, ok := l.Storage[id]; !ok {
		return "", ErrLinkNotFound
	} else {
		return val, nil

	}
}
