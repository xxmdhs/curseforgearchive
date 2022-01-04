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
	ID                 int     `db:"id"`
	Name               string  `db:"name"`
	Downloads          float64 `db:"downloadCount"`
	GameID             int     `db:"gameId"`
	GameCategoryID     int     `db:"gameCategoryId"`
	GamePopularityRank int     `db:"gamePopularityRank"`
	DateCreated        int64   `db:"dateCreated"`
	DateReleased       int64   `db:"dateReleased"`
	DateModified       int64   `db:"dateModified"`
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

func (s *Sqlite) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	return nil
}

func (s *Sqlite) Search(keyword string, limit, gameId, offset, sectionId int, ORDERBy string) ([]ModInfo, error) {
	var mods []ModInfo
	var r *sqlx.Rows
	var err error
	if sectionId == 0 {
		r, err = s.db.NamedQuery(`SELECT * FROM modinfo WHERE name LIKE :keyword AND gameId == :gameID ORDER BY "`+ORDERBy+`" DESC LIMIT :limit OFFSET :offset`,
			map[string]interface{}{"keyword": "%" + keyword + "%", "limit": limit, "offset": offset, "gameID": gameId, "sectionId": sectionId})

	} else {
		r, err = s.db.NamedQuery(`SELECT * FROM modinfo WHERE name LIKE :keyword AND gameId == :gameID AND gameCategoryId == :sectionId  ORDER BY "`+ORDERBy+`" DESC LIMIT :limit OFFSET :offset`,
			map[string]interface{}{"keyword": "%" + keyword + "%", "limit": limit, "offset": offset, "gameID": gameId, "sectionId": sectionId})
	}

	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}
	defer r.Close()
	for r.Next() {
		var m ModInfo
		err := r.StructScan(&m)
		if err != nil {
			return nil, fmt.Errorf("Search: %w", err)
		}
		mods = append(mods, m)
	}
	err = r.Err()
	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}
	return mods, nil
}
