package shortener

import (
	"context"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type UsecaseShortener struct {
	storage   repository.ShortenerRepo
	appConfig config.Config
}

func NewUsecase(localstorage repository.Storage, cfg config.Config) *UsecaseShortener {
	return &UsecaseShortener{
		storage:   &localstorage,
		appConfig: cfg,
	}
}

func (u *UsecaseShortener) CreateShortLink(ctx context.Context, longLink string) (link string) {

	switch {
	case u.appConfig.FileStoragePath == "":
		id, err := u.storage.CheckLinkInMemoryStorage(ctx, longLink)
		if err != nil {
			id = u.storage.SaveLinkInMemoryStorage(ctx, longLink)
			shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
			return shortLink
		} else {
			shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
			return shortLink
		}

	case u.appConfig.FileStoragePath != "":
		id, err := u.storage.CheckLinkInFileStorage(ctx, longLink)
		if err != nil {
			id = u.storage.SaveLinkInFileStorage(ctx, longLink)
			shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
			return shortLink
		} else {
			shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
			return shortLink
		}
	}
	return
}

func (u *UsecaseShortener) GetLinkByID(ctx context.Context, id int) (string, error) {

	if u.appConfig.FileStoragePath == "" {
		origLink, err := u.storage.GetLinkFromInMemoryStorage(ctx, id)
		if err == repository.ErrLinkNotFound {
			return "", err
		}
		return origLink, nil
	} else {
		origLink, err := u.storage.GetLinkFromInFileStorage(ctx, id)
		if err == repository.ErrLinkNotFound {
			return "", err
		}
		return origLink.LongLink, nil
	}
}
