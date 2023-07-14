package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/rltsv/urlcutter/internal/app/shortener/entity"
)

type PsqlStorage struct {
	db *pgx.Conn
}

func NewPsqlStorage(db *pgx.Conn) *PsqlStorage {
	return &PsqlStorage{
		db: db,
	}
}

func (ps *PsqlStorage) Ping(ctx context.Context) error {
	return ps.db.Ping(ctx)
}

func (ps *PsqlStorage) SaveLink(ctx context.Context, dto entity.Link) (userid, shorturl string, err error) {

	query := `INSERT INTO links (link_id, user_id, original_url, short_url, correlation_id) VALUES ($1,$2,$3,$4,$5);`

	_, err = ps.CheckOriginalURLInMemory(ctx, dto)
	if err != nil && err == ErrLinkNotFound {
		_, err = ps.db.Exec(ctx, query, dto.LinkID, dto.UserID, dto.OriginalURL, dto.ShortURL, dto.CorrelationID)
		if err != nil {
			return "", "", err
		}

		return dto.UserID, dto.ShortURL, nil
	}
	return "", "", ErrLinkAlreadyExist
}

func (ps *PsqlStorage) GetLink(ctx context.Context, dto entity.Link) (longurl string, err error) {

	query := `SELECT original_url FROM links WHERE link_id = $1 AND user_id = $2`

	row := ps.db.QueryRow(ctx, query, dto.LinkID, dto.UserID)

	err = row.Scan(&longurl)
	if err != nil && err == pgx.ErrNoRows {
		return "", ErrLinkNotFound
	}
	return longurl, nil
}

func (ps *PsqlStorage) GetLinksByUser(ctx context.Context, dto entity.Link) (links []entity.SendLinkDTO, err error) {
	query := `SELECT short_url, original_url FROM links WHERE user_id = $1`

	rows, err := ps.db.Query(ctx, query, dto.UserID)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		var link entity.SendLinkDTO
		err = rows.Scan(&link.ShortURL, &link.OriginalURL)
		if err != nil {
			return nil, err
		}
		err = rows.Err()
		if err != nil {
			log.Print(err)
		}

		links = append(links, link)
	}

	return links, nil
}

func (ps *PsqlStorage) CheckOriginalURLInMemory(ctx context.Context, dto entity.Link) (longurl string, err error) {
	query := `SELECT short_url FROM links WHERE original_url = $1 AND user_id = $2`

	row := ps.db.QueryRow(ctx, query, dto.OriginalURL, dto.UserID)

	err = row.Scan(&longurl)
	if err != nil && err == pgx.ErrNoRows {
		return "", ErrLinkNotFound
	}
	return longurl, nil
}
