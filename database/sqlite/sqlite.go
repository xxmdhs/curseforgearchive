package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	db *sqlx.DB
}

func NewSqlite(path string) (*Sqlite, error) {
	db, err := sqlx.Open("sqlite3", path+"?_txlock=IMMEDIATE&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("NewDb: %w", err)
	}
	s := Sqlite{}
	s.db = db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS modinfo(
		id INT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		downloadCount INT NOT NULL,
		gameId INT NOT NULL,
		gameCategoryId INT NOT NULL,
		gamePopularityRank INT NOT NULL,
		dateCreated INT NOT NULL,
		dateReleased INT NOT NULL,
		dateModified INT NOT NULL
	);`)
	if err != nil {
		return nil, fmt.Errorf("NewDb: %w", err)
	}

	return &s, nil
}

type ModInfo struct {
	ID                 int    `db:"id"`
	Name               string `db:"name"`
	Downloads          int    `db:"downloadCount"`
	GameID             int    `db:"gameId"`
	GameCategoryID     int    `db:"gameCategoryId"`
	GamePopularityRank int    `db:"gamePopularityRank"`
	DateCreated        int64  `db:"dateCreated"`
	DateReleased       int64  `db:"dateReleased"`
	DateModified       int64  `db:"dateModified"`
}

func (s *Sqlite) SetModInfo(m ModInfo) error {
	_, err := s.db.NamedExec(`INSERT OR REPLACE INTO modinfo
		(id, name, downloadCount, gameId, gameCategoryId, gamePopularityRank, dateCreated, dateReleased, dateModified)
		VALUES (:id, :name, :downloadCount, :gameId, :gameCategoryId, :gamePopularityRank, :dateCreated, :dateReleased, :dateModified)`, m)
	if err != nil {
		return fmt.Errorf("SetModInfo: %w", err)
	}
	return nil
}
