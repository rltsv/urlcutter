package shortener

import (
	"context"
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
	switch {
	case u.appConfig.FileStoragePath == "":
		return u.MemoryStorage.SaveLinkInMemoryStorage(ctx, link)
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

func (u *UsecaseShortener) GetLinkByUserID(ctx context.Context, dto entity.GetLinkDTO) (longurl string, err error) {
	LinkDTO := entity.GetLink(dto)

	switch {
	case u.appConfig.FileStoragePath == "":
		longurl, err = u.MemoryStorage.GetLinkFromInMemoryStorage(ctx, LinkDTO)
		if err != nil {
			return longurl, err
		}
		return longurl, err
	}
	if u.appConfig.FileStoragePath == "" {

	} else {

	}
	return longurl, err
}
