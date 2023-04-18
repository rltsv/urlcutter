package repository

import (
	"context"
	"github.com/rltsv/urlcutter/internal/app/config"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
	"os"
	"sync"
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

func (s *FileStorage) SaveLinkInFileStorage(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {
	//file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	//if err != nil {
	//	return userid, shorturl, err
	//}
	//defer file.Close()
	//
	//data, err := json.Marshal(value)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//writer := bufio.NewWriter(file)
	//
	//_, err = writer.Write(data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = writer.WriteByte('\n')
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = writer.Flush()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	return userid, shorturl, err
}

func (s *FileStorage) GetLinkFromInFileStorage(ctx context.Context, id int) (link ValueToFile, err error) {
	//	file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	//	if err != nil {
	//		log.Print("error while open file for storage")
	//	}
	//	defer file.Close()
	//
	//	scanner := bufio.NewScanner(file)
	//
	//	link = ValueToFile{}
	//
	//	for scanner.Scan() {
	//		err = json.Unmarshal(scanner.Bytes(), &link)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		if link.ID == id {
	//			return link, nil
	//		} else {
	//			continue
	//		}
	//	}
	//	return ValueToFile{}, ErrLinkNotFound
	//}
	//
	//func (s *FileStorage) CheckLinkInFileStorage(ctx context.Context, longLink string) (id int, err error) {
	//	file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	//	if err != nil {
	//		log.Print("error while open file for storage")
	//	}
	//	defer file.Close()
	//
	//	scanner := bufio.NewScanner(file)
	//
	//	link := ValueToFile{}
	//
	//	for scanner.Scan() {
	//
	//		err = json.Unmarshal(scanner.Bytes(), &link)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		if link.LongLink == longLink {
	//			return link.ID, nil
	//		}
	//	}
	return ValueToFile{}, ErrLinkNotFound
}

func (s *FileStorage) CheckMaxID() (id int) {
	//
	//file, err := os.OpenFile(s.appConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	//if err != nil {
	//	log.Print("error while open file for storage")
	//}
	//defer file.Close()
	//
	//scanner := bufio.NewScanner(file)
	//
	//link := ValueToFile{}
	//
	//var maxID = 0
	//
	//for scanner.Scan() {
	//	err = json.Unmarshal(scanner.Bytes(), &link)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	if link.ID > maxID {
	//		maxID = link.ID
	//	} else {
	//		continue
	//	}
	//}
	return 0
}

type ValueToFile struct {
	ID       int    `json:"id"`
	LongLink string `json:"long_link"`
}

var values []ValueToFile
