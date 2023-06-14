package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

type FileStorage struct {
	File      *os.File
	Mux       *sync.RWMutex
	appConfig config.Config
}

func NewFileStorage(cfg config.Config) *FileStorage {
	return &FileStorage{
		File:      new(os.File),
		Mux:       new(sync.RWMutex),
		appConfig: cfg,
	}
}

func (s *FileStorage) SaveLink(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {
	file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return userid, shorturl, err
	}
	defer file.Close()

	if ok := s.checkLinkInByUser(file, dto); ok {
		return userid, shorturl, ErrLinkAlreadyExist
	}

	if err = json.NewEncoder(file).Encode(&dto); err != nil {
		return userid, shorturl, err
	} else {
		return dto.UserID, dto.ShortURL, nil
	}
}

func (s *FileStorage) GetLink(ctx context.Context, dto entity.Link) (longurl string, err error) {
	file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return longurl, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		link := entity.Link{}
		err = json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Fatal(err)
		}

		if link.LinkID == dto.LinkID && link.UserID == dto.UserID {
			longurl = link.OriginalURL
			break
		} else {
			continue
		}
	}

	if longurl == "" {
		return longurl, ErrLinkNotFound
	} else {
		return longurl, nil
	}
}

func (s *FileStorage) GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error) {
	file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		link := entity.Link{}
		err = json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Fatal(err)
		}
		if link.UserID == dto.UserID {
			links = append(links, entity.SendLinkDTO{
				ShortURL:    link.ShortURL,
				OriginalURL: link.OriginalURL,
			})
		} else {
			continue
		}
	}
	if len(links) == 0 {
		return nil, ErrLinkNotFound
	}
	return links, nil
}

func (s *FileStorage) checkLinkInByUser(file *os.File, dto entity.Link) bool {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		link := entity.Link{}
		err := json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Print("failed while decode line in struct while scan")
		}
		if link.OriginalURL == dto.OriginalURL && link.UserID == dto.UserID {
			return true
		}
	}
	return false
}

func (s *FileStorage) Ping(ctx context.Context) error {
	return errors.New("there is no db in this configuration")
}
