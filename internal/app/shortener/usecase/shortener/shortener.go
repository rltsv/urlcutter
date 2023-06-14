package shortener

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type UsecaseShortener struct {
	storage   repository.Repository
	appConfig config.Config
}

func NewUsecase(storage repository.Repository, cfg config.Config) *UsecaseShortener {
	return &UsecaseShortener{
		storage:   storage,
		appConfig: cfg,
	}
}

func (u *UsecaseShortener) CreateShortLink(ctx context.Context, dto entity.CreateLinkDTO) (userid, shorturl string, err error) {
	link := entity.NewLink(dto)
	linkToSave := CreateNewLink(u.appConfig.BaseURL, link)

	return u.storage.SaveLink(ctx, *linkToSave)
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
	return u.storage.Ping(ctx)
}

//func (u *UsecaseShortener) BatchShortener(ctx context.Context, dto []entity.CreateLinkDTO) error {
//	var link entity.Link
//	var linkResponse []entity.SendLinkDTO
//	for _, val := range dto {
//		link.CorrelationID = val.CorrelationID
//		link.OriginalURL = val.OriginalURL
//		link.UserID = val.UserID
//
//		userid, shorturl, err := u.storage.SaveLink(ctx, link)
//		if err != nil {
//			return err
//		}
//
//	}
//}

// CreateNewLink create new instance of link
func CreateNewLink(baseurl string, dto entity.Link) *entity.Link {

	encodedLinkID := hex.EncodeToString(GenerateLinkID())
	switch {
	case dto.UserID != "":
		return &entity.Link{
			LinkID:      encodedLinkID,
			UserID:      dto.UserID,
			OriginalURL: dto.OriginalURL,
			ShortURL:    fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
		}
	default:
		encodedUserID := hex.EncodeToString(GenerateUserID())
		return &entity.Link{
			LinkID:      encodedLinkID,
			UserID:      encodedUserID,
			OriginalURL: dto.OriginalURL,
			ShortURL:    fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
		}
	}
}

// GenerateLinkID generate LinkID
func GenerateLinkID() []byte {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("error while generateLinkID: %v\n", err)
		return nil
	}
	return b
}

// GenerateUserID generate userID
func GenerateUserID() []byte {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("error while generateUserID: %v\n", err)
		return nil
	}
	return b
}
