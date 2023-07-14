package shortener

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
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

func (u *UsecaseShortener) BatchShortener(ctx context.Context, dto []entity.CreateLinkDTO) ([]entity.SendLinkDTO, error) {
	var response []entity.SendLinkDTO
	for _, val := range dto {
		// bring it to the link type
		link := entity.NewLink(val)
		// creating a link entity with all the necessary fields
		linkToSave := CreateNewLink(u.appConfig.BaseURL, link)

		_, shorturl, err := u.storage.SaveLink(ctx, *linkToSave)
		switch err {
		case repository.ErrLinkAlreadyExist:
			continue
		case nil:
			response = append(response, entity.SendLinkDTO{
				CorrelationID: link.CorrelationID,
				ShortURL:      shorturl,
			})
		}
	}
	if response == nil {
		return nil, errors.New("received links already shorten")
	}
	return response, nil
}

// CreateNewLink create new instance of link
func CreateNewLink(baseurl string, dto entity.Link) *entity.Link {

	encodedLinkID := hex.EncodeToString(GenerateLinkID())
	switch {
	case dto.UserID != "":
		return &entity.Link{
			LinkID:        encodedLinkID,
			UserID:        dto.UserID,
			OriginalURL:   dto.OriginalURL,
			ShortURL:      fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
			CorrelationID: dto.CorrelationID,
		}
	default:
		encodedUserID := hex.EncodeToString(GenerateUserID())
		return &entity.Link{
			LinkID:        encodedLinkID,
			UserID:        encodedUserID,
			OriginalURL:   dto.OriginalURL,
			ShortURL:      fmt.Sprintf("%s/%s", baseurl, encodedLinkID),
			CorrelationID: dto.CorrelationID,
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
