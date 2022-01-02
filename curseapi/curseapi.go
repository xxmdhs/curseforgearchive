package curseapi

import (
	"fmt"
	"net/url"
	"strconv"
)

//From https://gaz492.github.io/TwitchAPI/

func Searchmod(key string, index, sectionId int) ([]byte, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/search?categoryId=0&gameId=432&index=` +
		strconv.Itoa(index) + `&pageSize=20&searchFilter=` + url.QueryEscape(key) + `&sectionId=` +
		strconv.Itoa(sectionId) + `&sort=5`
	b, err := httpget(aurl)
	if err != nil {
		return nil, fmt.Errorf("Searchmod: %w", err)
	}
	return b, nil
}

func Addonfiles(addonID string) ([]byte, error) {
	aurl := `https://addons-ecs.forgesvc.net/api/v2/addon/` + addonID + `/files`
	b, err := httpget(aurl)
	if err != nil {
		return nil, fmt.Errorf("Addonfiles: %w", err)
	}
	return b, nil
}
