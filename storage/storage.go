package storage

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"

	"anio/storage/migrations"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
)

const dbName = "db.sqlite3"

type Storage struct {
	db *bun.DB
}

func New(path string) (*Storage, error) {
	path = filepath.Join(path, dbName)
	dsn := fmt.Sprintf("file:%s?mode=rwc", path)

	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	isDebug := zerolog.GlobalLevel() == zerolog.DebugLevel

	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(isDebug),
		bundebug.WithVerbose(isDebug),
	))

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	st := &Storage{db: db}

	// st.db.NewCreateTable().Table("test_table").Exec(context.Background())

	if err = st.migrate(context.Background()); err != nil {
		return nil, fmt.Errorf("migrations error: %w", err)
	}

	return st, nil
}

func (st *Storage) migrate(ctx context.Context) error {
	migrator := migrate.NewMigrator(st.db, migrations.Migrations)
	if err := migrator.Init(context.Background()); err != nil {
		return fmt.Errorf("migrator init error: %w", err)
	}

	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		return fmt.Errorf("migrations apply error: %w", err)
	}

	if group.IsZero() {
		log.Debug().Msgf("there are no new migrations to apply")
		return nil
	}

	log.Debug().Msgf("migrated to %s", group)
	return nil
}
