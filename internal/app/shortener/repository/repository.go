package repository

import (
	"context"
	"errors"
)

var (
	ErrLinkNotFound    = errors.New("link not found")
	ErrRepositoryEmpty = errors.New("there are no any saved links")
)

type ShortenerRepo interface {
	SaveLinkInMemoryStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInMemoryStorage(ctx context.Context, id int) (longLink string, err error)
	CheckLinkInMemoryStorage(ctx context.Context, longLink string) (id int, err error)

	SaveLinkInFileStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInFileStorage(ctx context.Context, id int) (link ValueToFile, err error)
	CheckLinkInFileStorage(ctx context.Context, longLink string) (id int, err error)
}

