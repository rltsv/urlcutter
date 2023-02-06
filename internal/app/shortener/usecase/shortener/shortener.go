package shortener

import (
	"context"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type Usecase struct {
	repo repository.ShortenerRepo
}

func NewUsecase(shortenerRepo repository.ShortenerRepo) *Usecase {
	return &Usecase{repo: shortenerRepo}
}

func (u *Usecase) CreateShortLink(ctx context.Context, longLink string) (link string) {

	IDCount := u.repo.CreateLink(ctx, longLink)

	shortLink := fmt.Sprintf("%s%d", config.Cfg.BaseURL, IDCount)

	return shortLink
}

func (u *Usecase) GetLinkByID(ctx context.Context, id int) (origLink string, err error) {

	origLink, err = u.repo.GetLinkByID(ctx, id)
	if err == repository.ErrLinkNotFound {
		return "", err
	}

	return origLink, nil
}
