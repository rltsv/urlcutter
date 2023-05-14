package shortener

import (
	"context"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type UsecaseShortener struct {
	storage   repository.Repository
	db        repository.DBRepository
	appConfig config.Config
}

func NewUsecase(storage repository.Repository, db repository.DBRepository, cfg config.Config) *UsecaseShortener {
	return &UsecaseShortener{
		storage:   storage,
		db:        db,
		appConfig: cfg,
	}
}

func (u *UsecaseShortener) CreateShortLink(ctx context.Context, dto entity.CreateLinkDTO) (userid, shorturl string, err error) {
	link := entity.NewLink(dto)
	return u.storage.SaveLink(ctx, link)

}

func (u *UsecaseShortener) GetLinkByUserID(ctx context.Context, dto entity.GetLinkDTO) (longurl string, err error) {
	link := entity.GetLink(dto)
	return u.storage.GetLink(ctx, link)
}

func (u *UsecaseShortener) GetLinksByUser(ctx context.Context, dto entity.GetAllLinksDTO) (links []entity.SendLinkDTO, err error) {
	user := entity.GetAllLinks(dto)
	return u.storage.GetLinksByUser(ctx, user)
}

func (u *UsecaseShortener) Ping(ctx context.Context) error {
	return u.db.Ping(ctx)
}
