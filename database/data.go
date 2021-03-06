package database

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDB struct {
	Db *leveldb.DB
}

func NewLevelDB(path string) (*LevelDB, error) {
	l := &LevelDB{}
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("NewLevelDB: %w", err)
	}
	l.Db = db
	return l, nil
}

func (l *LevelDB) Close() error {
	err := l.Db.Close()
	if err != nil {
		return fmt.Errorf("LevelDB.Close: %w", err)
	}
	return nil
}

func (l *LevelDB) Get(key string) ([]byte, error) {
	val, err := l.Db.Get([]byte(key), nil)
	if err != nil {
		return nil, fmt.Errorf("LevelDB.Get: %w", err)
	}
	return val, nil
}

func (l *LevelDB) Put(key string, value []byte) error {
	err := l.Db.Put([]byte(key), value, nil)
	if err != nil {
		return fmt.Errorf("LevelDB.Put: %w", err)
	}
	return nil
}
