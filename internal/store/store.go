package store

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func sqliteDSN(path string) string {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	return path + separator + "_pragma=busy_timeout(5000)"
}

func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", sqliteDSN(path))
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database %s: %w", path, err)
	}

	database, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create sqlite migration driver: %w", err)
	}

	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "sqlite", database)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
