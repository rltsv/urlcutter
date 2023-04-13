package repository

import (
	"context"
	"errors"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

var (
	ErrLinkNotFound    = errors.New("link not found")
	ErrRepositoryEmpty = errors.New("there are no any saved links")
)

type MemoryRepository interface {
	SaveLinkInMemoryStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error)
	GetLinkFromInMemoryStorage(ctx context.Context, dto entity.Link) (longurl string, err error)
	CheckLinkInMemoryStorage(ctx context.Context, linkdto entity.Link) (id int, err error)
}

type FileRepository interface {
	SaveLinkInFileStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInFileStorage(ctx context.Context, id int) (link ValueToFile, err error)
	CheckLinkInFileStorage(ctx context.Context, longLink string) (id int, err error)
}
