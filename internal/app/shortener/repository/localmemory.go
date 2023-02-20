package repository

import (
	"context"
	"encoding/json"
	"github.com/rltsv/urlcutter/internal/app/config"
	"log"
	"os"
	"strings"
	"sync"
)

type Storage struct {
	InMemoryStorage map[int]string
	InFileStorage   *os.File
	IDCount         int
	Mux             *sync.RWMutex
	Decoder         *json.Decoder
	Encoder         *json.Encoder
}

func NewStorage() *Storage {
	file, err := os.OpenFile(config.Cfg.FileStoragePath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Print("error while open file for storage")
	}
	return &Storage{
		InMemoryStorage: make(map[int]string),
		InFileStorage:   file,
		IDCount:         0,
		Mux:             new(sync.RWMutex),
		Encoder:         json.NewEncoder(file),
		Decoder:         json.NewDecoder(file),
	}
}
func (l *Storage) SaveLinkInMemoryStorage(ctx context.Context, longLink string) (id int) {
	id = l.CheckLinkInMemoryStorage(longLink)
	if id != 0 {
		return id
	}
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

func (l *Storage) CheckLinkInMemoryStorage(longLink string) (id int) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()
	for key, val := range l.InMemoryStorage {
		if strings.EqualFold(val, longLink) {
			return key
		}
	}
	return 0
}
