package repository

import (
	"context"
	"log"
	"strings"
)

type ValueToFile struct {
	ID       int    `json:"id"`
	LongLink string `json:"long_link"`
}

var values []ValueToFile

func (l *Storage) SaveLinkInFileStorage(ctx context.Context, longLink string) (id int) {
	id = l.CheckLinkInFileStorage(longLink)
	if id != 0 {
		return id
	}

	l.Mux.Lock()
	defer l.Mux.Unlock()

	err := l.Decoder.Decode(&values)
	if err != nil {
		log.Print(ErrWhileEncode)
	}

	l.IDCount++
	value := ValueToFile{
		ID:       l.IDCount,
		LongLink: longLink,
	}
	values := append(values, value)

	err = l.Encoder.Encode(&values)
	if err != nil {
		log.Print(ErrWhileEncode)
	}

	return 0
}

func (l *Storage) GetLinkFromInFileStorage(ctx context.Context, id int) (longLink string, err error) {
	l.Mux.RLock()
	defer l.Mux.RUnlock()

	err = l.Decoder.Decode(&values)
	if err != nil {
		log.Print(ErrWhileDecode)
	}

	for _, val := range values {
		if val.ID == id {
			return val.LongLink, nil
		}
	}
	return "", nil
}

func (l *Storage) CheckLinkInFileStorage(longLink string) (id int) {

	l.Mux.RLock()
	defer l.Mux.RUnlock()

	err := l.Decoder.Decode(&values)
	if err != nil {
		log.Print(ErrWhileEncode)
	}

	for _, val := range values {
		if strings.EqualFold(val.LongLink, longLink) {
			return val.ID
		}
	}
	return 0
}
