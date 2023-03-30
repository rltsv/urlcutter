package repository

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
)

type ValueToFile struct {
	ID       int    `json:"id"`
	LongLink string `json:"long_link"`
}

var values []ValueToFile

func (l *Storage) SaveLinkInFileStorage(ctx context.Context, longLink string) (id int) {
	file, err := os.OpenFile(l.AppConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Print("error while open file for storage")
	}
	defer file.Close()

	maxID := l.CheckMaxID()
	l.IDCount = maxID

	l.IDCount++
	value := ValueToFile{
		ID:       l.IDCount,
		LongLink: longLink,
	}

	log.Print(value)

	data, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(file)

	_, err = writer.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.WriteByte('\n')
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}

	return l.IDCount
}

func (l *Storage) GetLinkFromInFileStorage(ctx context.Context, id int) (link ValueToFile, err error) {
	file, err := os.OpenFile(l.AppConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Print("error while open file for storage")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	link = ValueToFile{}

	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Fatal(err)
		}

		if link.ID == id {
			return link, nil
		} else {
			continue
		}
	}
	return ValueToFile{}, ErrLinkNotFound
}

func (l *Storage) CheckLinkInFileStorage(ctx context.Context, longLink string) (id int, err error) {
	file, err := os.OpenFile(l.AppConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Print("error while open file for storage")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	link := ValueToFile{}

	for scanner.Scan() {

		err = json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Fatal(err)
		}

		if link.LongLink == longLink {
			return link.ID, nil
		}
	}
	return 0, ErrLinkNotFound
}

func (l *Storage) CheckMaxID() (id int) {

	file, err := os.OpenFile(l.AppConfig.FileStoragePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		log.Print("error while open file for storage")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	link := ValueToFile{}

	var maxID = 0

	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &link)
		if err != nil {
			log.Fatal(err)
		}

		if link.ID > maxID {
			maxID = link.ID
		} else {
			continue
		}
	}
	return maxID
}
