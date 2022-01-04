package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/xxmdhs/curseforgearchive/database"
	"github.com/xxmdhs/curseforgearchive/database/sqlite"
)

func NewServer(addr string, db *database.LevelDB, sqlite *sqlite.Sqlite) error {
	r := httprouter.New()
	server := newServer(db, sqlite)
	defer server.Close()

	r.GET("/api/v2/addon/:addonID", server.mod("modinfo-"))
	r.GET("/api/v2/addon/:addonID/files", server.mod("modfiles-"))
	r.GET("/api/v2/addon/:addonID/file/:fileId", server.modfile(false))
	r.GET("/api/v2/addon/:addonID/file/:fileId/download-url", server.modfile(true))

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
