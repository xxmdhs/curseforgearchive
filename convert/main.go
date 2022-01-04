package main

import (
	"encoding/json"
	"time"

	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/xxmdhs/curseforgearchive/database"
	"github.com/xxmdhs/curseforgearchive/database/sqlite"
)

func main() {
	db, err := database.NewLevelDB("./data")
	e(err)

	d, err := sqlite.NewSqlite("./data.db")
	e(err)

	i := db.Db.NewIterator(util.BytesPrefix([]byte("modinfo-")), nil)
	defer i.Release()
	for i.Next() {
		info := modinfo{}
		e(json.Unmarshal(i.Value(), &info))

		DateCreated, err := time.Parse(time.RFC3339Nano, info.DateCreated)
		e(err)
		DateReleased, err := time.Parse(time.RFC3339Nano, info.DateReleased)
		e(err)
		DateModified, err := time.Parse(time.RFC3339Nano, info.DateModified)
		e(err)

		m := sqlite.ModInfo{
			ID:                 info.ID,
			Name:               info.Name,
			Downloads:          info.DownloadCount,
			GameID:             info.GameID,
			GameCategoryID:     info.CategorySection.GameCategoryID,
			GamePopularityRank: info.GamePopularityRank,
			DateCreated:        DateCreated.Unix(),
			DateReleased:       DateReleased.Unix(),
			DateModified:       DateModified.Unix(),
		}
		d.SetModInfo(m)
	}
	e(i.Error())
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

type modinfo struct {
	GamePopularityRank int                    `json:"gamePopularityRank"`
	CategorySection    modinfoCategorySection `json:"categorySection"`
	DateCreated        string                 `json:"dateCreated"`
	DateModified       string                 `json:"dateModified"`
	DateReleased       string                 `json:"dateReleased"`
	DownloadCount      int                    `json:"downloadCount"`
	GameID             int                    `json:"gameId"`
	ID                 int                    `json:"id"`
	Name               string                 `json:"name"`
}

type modinfoCategorySection struct {
	GameCategoryID          int    `json:"gameCategoryId"`
	GameID                  int    `json:"gameId"`
	ID                      int    `json:"id"`
	InitialInclusionPattern string `json:"initialInclusionPattern"`
	Name                    string `json:"name"`
	PackageType             int    `json:"packageType"`
	Path                    string `json:"path"`
}
