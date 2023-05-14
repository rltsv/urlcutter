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

type Repository interface {
	SaveLink(ctx context.Context, dto entity.Link) (userid, shorturl string, err error)
	GetLink(ctx context.Context, dto entity.Link) (longurl string, err error)
	GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error)
}

type DBRepository interface {
	Ping(ctx context.Context) error
}
