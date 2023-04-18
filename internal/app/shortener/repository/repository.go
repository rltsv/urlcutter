package repository

import (
	"context"
	"errors"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

var (
	ErrRepositoryEmpty  = errors.New("there are no any saved links")
	ErrLinkAlreadyExist = errors.New("this link already shortened")
	ErrUnknownLink      = errors.New("unknown link")
	ErrLinkNotFound     = errors.New("link not found")
	ErrUserIsNotFound   = errors.New("there is no any user in memory with this id")
)

type MemoryRepository interface {
	SaveLinkInMemoryStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error)
	GetLinkFromInMemoryStorage(ctx context.Context, dto entity.Link) (longurl string, err error)
	GetLinksByUser(ctx context.Context, dto entity.Link) []entity.SendLinkDTO
}

type FileRepository interface {
	SaveLinkInFileStorage(ctx context.Context, longLink string) (id int)
	GetLinkFromInFileStorage(ctx context.Context, id int) (link ValueToFile, err error)
	CheckLinkInFileStorage(ctx context.Context, longLink string) (id int, err error)
}
