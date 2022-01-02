package curseapi

import (
	"encoding/json"
	"fmt"
)

type Modinfo struct {
	Name                   string             `json:"name"`
	ID                     int                `json:"id"`
	GameVersionLatestFiles []GameVersionFiles `json:"gameVersionLatestFiles"`
	WebsiteUrl             string             `json:"websiteUrl"`
	Summary                string             `json:"summary"`
}

type GameVersionFiles struct {
	GameVersion     string `json:"gameVersion"`
	ProjectFileId   int    `json:"projectFileId"`
	ProjectFileName string `json:"projectFileName"`
	json.RawMessage
}

func json2Modinfo(jsonbyte []byte) ([]Modinfo, error) {
	m := make([]Modinfo, 0, 5)
	err := json.Unmarshal(jsonbyte, &m)
	if err != nil {
		return nil, fmt.Errorf("json2Modinfo: %w", err)
	}
	return m, nil
}

type Files struct {
	ID           int            `json:"id"`
	FileName     string         `json:"fileName"`
	FileDate     string         `json:"fileDate"`
	Dependencies []Dependencies `json:"dependencies"`
	DownloadUrl  string         `json:"downloadUrl"`
	GameVersion  []string       `json:"gameVersion"`
	ReleaseType  int            `json:"releaseType"`
	json.RawMessage
}

type Dependencies struct {
	AddonId int `json:"addonId"`
	Type    int `json:"type"`
	json.RawMessage
}
