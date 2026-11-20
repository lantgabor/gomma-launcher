package store

import (
	"database/sql"
	"errors"
	"fmt"
)

func (s *Store) GetRevision(url string) (string, bool, error) {
	var revision string
	err := s.db.QueryRow(`SELECT revision FROM git_revisions WHERE url = ?`, url).Scan(&revision)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("failed to get revision %q: %w", url, err)
	}
	return revision, true, nil
}

func (s *Store) SetRevision(url, revision string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`INSERT INTO git_revisions (url, revision) VALUES (?, ?)
		 ON CONFLICT(url) DO UPDATE SET revision = excluded.revision`,
		url, revision,
	)
	if err != nil {
		return fmt.Errorf("failed to set revision %q: %w", url, err)
	}
	return nil
}
