package store

import (
	"database/sql"
	"errors"
	"fmt"
)

type Mod struct {
	Name            string
	ArchiveFilename string
	DownloadedMD5   string
	InstalledMD5    string
}

func (s *Store) GetMod(name string) (*Mod, bool, error) {
	var m Mod
	err := s.db.QueryRow(
		`SELECT name, archive_filename, downloaded_md5, installed_md5
		 FROM mods WHERE name = ?`,
		name,
	).Scan(&m.Name, &m.ArchiveFilename, &m.DownloadedMD5, &m.InstalledMD5)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get mod %q: %w", name, err)
	}
	return &m, true, nil
}

func (s *Store) UpsertMod(m Mod) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`INSERT INTO mods (name, archive_filename, downloaded_md5, installed_md5)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET
		     archive_filename = excluded.archive_filename,
		     downloaded_md5   = excluded.downloaded_md5,
		     installed_md5    = excluded.installed_md5`,
		m.Name, m.ArchiveFilename, m.DownloadedMD5, m.InstalledMD5,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert mod %q: %w", m.Name, err)
	}
	return nil
}

func (s *Store) SetDownloaded(name, filename, md5 string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`INSERT INTO mods (name, archive_filename, downloaded_md5)
		 VALUES (?, ?, ?)
		 ON CONFLICT(name) DO UPDATE SET
		     archive_filename = excluded.archive_filename,
		     downloaded_md5   = excluded.downloaded_md5`,
		name, filename, md5,
	)
	if err != nil {
		return fmt.Errorf("failed to set downloaded mod %q: %w", name, err)
	}
	return nil
}

func (s *Store) SetInstalled(name, md5 string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`INSERT INTO mods (name, installed_md5)
		 VALUES (?, ?)
		 ON CONFLICT(name) DO UPDATE SET installed_md5 = excluded.installed_md5`,
		name, md5,
	)
	if err != nil {
		return fmt.Errorf("failed to set installed mod %q: %w", name, err)
	}
	return nil
}
