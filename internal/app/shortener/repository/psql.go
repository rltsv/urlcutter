package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PsqlStorage struct {
	db *pgx.Conn
}

func NewPsqlStorage(db *pgx.Conn) *PsqlStorage {
	return &PsqlStorage{
		db: db,
	}
}

func (db *PsqlStorage) Ping(ctx context.Context) error {
	return db.db.Ping(ctx)
}
