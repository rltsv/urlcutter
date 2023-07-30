package repository

import (
	"context"
	"errors"

	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

var (
	ErrLinkAlreadyExist = errors.New("this link already shortened")
	ErrLinkNotFound     = errors.New("link not found")
	ErrUserIsNotFound   = errors.New("there is no any user in memory with this id")
)

type Repository interface {
	SaveLink(ctx context.Context, dto entity.Link) (userid, shorturl string, err error)
	GetLink(ctx context.Context, dto entity.Link) (longurl string, err error)
	GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error)
	Ping(ctx context.Context) error
}
