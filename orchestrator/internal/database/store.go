package database

import (
	"database/sql"
	"errors"
	"time"
	"log"
	_ "github.com/lib/pq"	
	"github.com/souls-syntax/Templates/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db:db}, nil

}


func (s *Store) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS queries (
		hash TEXT PRIMARY KEY,
		query_text TEXT NOT NULL,
		verdict TEXT,
		confidence FLOAT,
		decider TEXT,
		source TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		verified_by_human BOOLEAN DEFAULT FALSE,
		human_notes	TEXT,
		hit_count	INT DEFAULT 1
	);
	`
	_, err := s.db.Exec(query)
	return err
}

// 1. The writer

func (s *Store) SaveDecision(d models.Decision, source string) error {
	query := `
	INSERT INTO queries (hash, query_text, verdict, confidence, decider, source, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (hash) DO UPDATE
	SET verdict = EXCLUDED.verdict, confidence = EXCLUDED.confidence;
	`
	_, err := s.db.Exec(query, d.QueryHash, d.QueryText, d.Verdict, d.Confidence, d.Decider, source, time.Now())
	return err
}

// 2. The reader

func (s *Store) GetDecision(hash string) (models.Decision, error) {
	query := `SELECT query_text, verdict, confidence, decider FROM queries WHERE hash = $1`

	row := s.db.QueryRow(query, hash)

	var d models.Decision
	
	err := row.Scan(&d.QueryText, &d.Verdict, &d.Confidence, &d.Decider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Decision{}, errors.New("Not Found")
		}
		return models.Decision{}, err
	}

	updateQuery := `
		UPDATE queries
		SET hit_count = hit_count + 1
		WHERE hash = $1
	`
	if _, err := s.db.Exec(updateQuery, hash); err != nil {
		log.Printf("warn: failed to increment hit_count for hash %s: %v", hash, err)
	}
	return d, nil
}
