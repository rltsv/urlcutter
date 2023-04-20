package shortener

import (
	"context"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type UsecaseShortener struct {
	memoryStorage repository.MemoryRepository
	fileStorage   repository.FileRepository
	appConfig     config.Config
}

func NewUsecase(memorystorage repository.MemoryRepository, filestorage repository.FileRepository, cfg config.Config) *UsecaseShortener {
	return &UsecaseShortener{
		memoryStorage: memorystorage,
		fileStorage:   filestorage,
		appConfig:     cfg,
	}
}

func (u *UsecaseShortener) CreateShortLink(ctx context.Context, dto entity.CreateLinkDTO) (userid, shorturl string, err error) {
	link := entity.NewLink(dto)
	switch {
	case u.appConfig.FileStoragePath == "":
		return u.memoryStorage.SaveLinkInMemoryStorage(ctx, link)
	case u.appConfig.FileStoragePath != "":
		return u.fileStorage.SaveLinkInFileStorage(ctx, link)
	}
	return
}

func (u *UsecaseShortener) GetLinkByUserID(ctx context.Context, dto entity.GetLinkDTO) (longurl string, err error) {
	link := entity.GetLink(dto)

	switch {
	case u.appConfig.FileStoragePath == "":
		longurl, err = u.memoryStorage.GetLinkFromInMemoryStorage(ctx, link)
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

func (u *UsecaseShortener) GetLinksByUser(ctx context.Context, dto entity.GetAllLinksDTO) (links []entity.SendLinkDTO, err error) {
	user := entity.GetAllLinks(dto)

	switch {
	case u.appConfig.FileStoragePath == "":
		return u.memoryStorage.GetLinksByUser(ctx, user)
	case u.appConfig.FileStoragePath != "":
		//sdasd
	}
	return links, nil
}
