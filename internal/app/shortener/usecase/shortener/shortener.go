package shortener

import (
	"context"
	"errors"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type UsecaseShortener struct {
	MemoryStorage repository.MemoryStorage
	FileStorage   repository.FileStorage
	appConfig     config.Config
}

func NewUsecase(memorystorage repository.MemoryStorage, filestorage repository.FileStorage, cfg config.Config) *UsecaseShortener {
	return &UsecaseShortener{
		MemoryStorage: memorystorage,
		FileStorage:   filestorage,
		appConfig:     cfg,
	}
}

func (u *UsecaseShortener) CreateShortLink(ctx context.Context, dto entity.CreateLinkDTO) (userid, shorturl string, err error) {
	link := entity.NewLink(dto)
	if link.LongURL == "" {
		return "", "", errors.New("there is nothing to shorten")
	}

	switch {
	case u.appConfig.FileStoragePath == "":
		userid, shorturl, err = u.MemoryStorage.SaveLinkInMemoryStorage(ctx, link)
		if err != nil && err == repository.ErrLinkAlreadyExist {
			return "", "", err
		}

		return userid, shorturl, nil
	case u.appConfig.FileStoragePath != "":
		//id, err := u.storage.CheckLinkInFileStorage(ctx, longLink)
		//if err != nil {
		//	id = u.storage.SaveLinkInFileStorage(ctx, longLink)
		//	shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
		//	return shortLink
		//} else {
		//	shortLink := fmt.Sprintf("%s/%d", u.appConfig.BaseURL, id)
		//	return shortLink
		//}
	}
	return
}

func (u *UsecaseShortener) GetLinkByID(ctx context.Context, id int) (string, error) {

	if u.appConfig.FileStoragePath == "" {

	} else {

	}
	return "", nil
}
