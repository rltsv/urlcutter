package shortener

import (
	"context"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type Usecase struct {
	inMemoryStorage repository.ShortenerRepo
}

func NewUsecase(localstorage repository.Storage) *Usecase {
	return &Usecase{
		inMemoryStorage: &localstorage,
	}
}

func (u *Usecase) CreateShortLink(ctx context.Context, longLink string) (link string) {

	if config.Cfg.FileStoragePath == "" {
		IDCount := u.inMemoryStorage.SaveLinkInMemoryStorage(ctx, longLink)
		shortLink := fmt.Sprintf("%s/%d", config.Cfg.BaseURL, IDCount)
		return shortLink
	} else {
		IDCount := u.inMemoryStorage.SaveLinkInFileStorage(ctx, longLink)
		shortLink := fmt.Sprintf("%s/%d", config.Cfg.BaseURL, IDCount)
		return shortLink
	}
}

func (u *Usecase) GetLinkByID(ctx context.Context, id int) (origLink string, err error) {

	if config.Cfg.FileStoragePath == "" {
		origLink, err = u.inMemoryStorage.GetLinkFromInMemoryStorage(ctx, id)
		if err == repository.ErrLinkNotFound {
			return "", err
		}
		return origLink, nil
	} else {
		origLink, err = u.inMemoryStorage.GetLinkFromInFileStorage(ctx, id)
		if err == repository.ErrLinkNotFound {
			return "", err
		}
		return origLink, nil
	}
}
