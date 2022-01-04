package sqlite

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSqlite_Search(t *testing.T) {
	db, err := NewSqlite("./data.db")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	l, err := db.Search("jei", 20, 432, 0, 0, "dateCreated")
	if err != nil {
		t.Error(err)
	}
	b, _ := json.Marshal(l)
	fmt.Println(string(b))
}
