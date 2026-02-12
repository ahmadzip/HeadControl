package store

import (
	"database/sql"
	"headcontrol/internal/model"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	db *sql.DB
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			base_url TEXT NOT NULL,
			api_key TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *Store) GetSettings() (*model.Settings, error) {
	var st model.Settings
	err := s.db.QueryRow(
		"SELECT id, base_url, api_key, created_at, updated_at FROM settings ORDER BY id DESC LIMIT 1",
	).Scan(&st.ID, &st.BaseURL, &st.APIKey, &st.CreatedAt, &st.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *Store) SaveSettings(baseURL, apiKey string) error {
	existing, err := s.GetSettings()
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)

	if existing != nil {
		_, err = s.db.Exec(
			"UPDATE settings SET base_url = ?, api_key = ?, updated_at = ? WHERE id = ?",
			baseURL, apiKey, now, existing.ID,
		)
	} else {
		_, err = s.db.Exec(
			"INSERT INTO settings (base_url, api_key, created_at, updated_at) VALUES (?, ?, ?, ?)",
			baseURL, apiKey, now, now,
		)
	}
	return err
}

func (s *Store) HasSettings() bool {
	st, err := s.GetSettings()
	return err == nil && st != nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
