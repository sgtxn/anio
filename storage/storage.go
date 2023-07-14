package storage

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type Storage struct {
	db *bun.DB
}

func New(path string) (*Storage, error) {
	dsn := fmt.Sprintf("file::%s:?mode=rwc", path)
	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	return &Storage{db: db}, nil
}
