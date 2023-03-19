package repository

import (
	"context"
	"strings"
	"sync"
)

type Storage struct {
	InMemoryStorage map[int]string
	IDCount         int
	Mux             *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		InMemoryStorage: make(map[int]string),
		IDCount:         0,
		Mux:             new(sync.RWMutex),
	}
}
func (l *Storage) SaveLinkInMemoryStorage(ctx context.Context, longLink string) (id int) {
	l.Mux.Lock()
	defer l.Mux.Unlock()
	l.IDCount++
	l.InMemoryStorage[l.IDCount] = longLink
	return l.IDCount
}

func (l *Storage) GetLinkFromInMemoryStorage(ctx context.Context, id int) (longLink string, err error) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()
	if val, ok := l.InMemoryStorage[id]; !ok {
		return "", ErrLinkNotFound
	} else {
		return val, nil
	}
}

func (l *Storage) CheckLinkInMemoryStorage(ctx context.Context, longLink string) (id int, err error) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()
	for key, val := range l.InMemoryStorage {
		if strings.EqualFold(val, longLink) {
			return key, nil
		}
	}
	return 0, ErrLinkNotFound
}
