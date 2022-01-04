package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xxmdhs/curseforgearchive/database"
	"github.com/xxmdhs/curseforgearchive/database/sqlite"
)

func NewServer(addr string, db *database.LevelDB, sqlite *sqlite.Sqlite) error {
	r := httprouter.New()
	server := newServer(db, sqlite)
	defer server.Close()

	r.GET("/api/v2/addon/:addonID", server.mod("modinfo-"))
	r.GET("/api/v2/addon/:addonID/files", server.mod("modfiles-"))

	s := http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		return fmt.Errorf("NewServer: %w", err)
	}
	return nil
}

type server struct {
	db     *database.LevelDB
	sqlite *sqlite.Sqlite
}

func newServer(db *database.LevelDB, sqlite *sqlite.Sqlite) *server {
	return &server{
		db:     db,
		sqlite: sqlite,
	}
}

func (s *server) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	err = s.sqlite.Close()
	if err != nil {
		return fmt.Errorf("Close: %w", err)
	}
	return nil
}

func (s *server) mod(prefix string) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		addonID := p.ByName("addonID")
		if addonID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if addonID == "search" {
			s.search(w, r, p)
			return
		}

		b, err := s.db.Get(prefix + addonID)
		if err != nil {
			if !errors.Is(err, leveldb.ErrNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
	}
}

func (s *server) search(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	gameId, err := getINT(r, "gameId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	pageSize, err := getINT(r, "pageSize")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	sectionId, err := getINT(r, "sectionId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	sort, err := getINT(r, "sort")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	index, err := getINT(r, "index")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	keyword := r.FormValue("searchFilter")

	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	orderby := ""
	switch sort {
	case 0:
		orderby = "gamePopularityRank"
	case 1:
		orderby = "gamePopularityRank"
	case 2:
		orderby = "dateModified"
	case 3:
		orderby = "name"
	case 4:
		orderby = "gamePopularityRank"
	case 5:
		orderby = "downloadCount"
	default:
		orderby = "gamePopularityRank"
	}

	l, err := s.sqlite.Search(keyword, pageSize, gameId, index, sectionId, orderby)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	ls := []interface{}{}

	for _, v := range l {
		b, err := s.db.Get("modinfo-" + strconv.Itoa(v.ID))
		if err != nil {
			if !errors.Is(err, leveldb.ErrNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			continue
		}
		var in interface{}
		err = json.Unmarshal(b, &in)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		ls = append(ls, in)
	}

	b, err := json.Marshal(ls)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(b)
}

func getINT(r *http.Request, key string) (int, error) {
	v := r.FormValue(key)
	if v == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("getV: %w", err)
	}
	return i, nil
}
