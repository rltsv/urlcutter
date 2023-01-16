package shortener

import (
	"context"
	"fmt"
	"github.com/rltsv/urlcutter/internal/app/shortener/repository"
)

type Usecase struct {
	repo repository.ShortenerRepo
}

func NewUsecase(shortenerRepo repository.ShortenerRepo) *Usecase {
	return &Usecase{repo: shortenerRepo}
}

func (u *Usecase) CreateShortLink(ctx context.Context, longLink string) (link string, err error) {

	IDCount := u.repo.SaveLink(ctx, longLink)

	shortLink := fmt.Sprintf("http://localhost:8080/%d", IDCount)

	return shortLink, nil
}

func (u *Usecase) GetLinkByID(ctx context.Context, id int) (origLink string, err error) {

	origLink, err = u.repo.GetLink(ctx, id)
	if err == repository.ErrStorageIsEmpty {
		return "", err
	} else if err == repository.ErrLinkNotFound {
		return "", err
	}

	return origLink, nil
}
