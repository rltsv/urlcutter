package repository

import (
	"context"
	"errors"
)

var (
	ErrLinkNotFound = errors.New("link not found")
	ErrWhileDecode  = errors.New("error occurred while decoding in file")
	ErrWhileEncode  = errors.New("error occurred while encode in file")
)

type ShortenerRepo interface {
	SaveLinkInMemoryStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInMemoryStorage(ctx context.Context, id int) (longLink string, err error)
	SaveLinkInFileStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInFileStorage(ctx context.Context, id int) (longLink string, err error)
}
