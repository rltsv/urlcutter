package repository

import (
	"context"
	"errors"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"os"
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
	GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error)
}

type FileRepository interface {
	SaveLinkInFileStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error)
	GetLinkFromFileStorage(ctx context.Context, id int) (err error)
	checkLinkInByUser(file *os.File, dto entity.Link) bool
}
