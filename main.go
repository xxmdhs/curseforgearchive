package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xxmdhs/curseforgearchive/curseapi"
	"github.com/xxmdhs/curseforgearchive/database"
)

func main() {
	b, err := os.ReadFile("config.json")
	e(err)
	c := config{}
	e(json.Unmarshal(b, &c))

	db, err := database.NewLevelDB("./data")
	e(err)
	defer db.Close()

	start := 0
	if !c.Reacquire {
		b, err := db.Get("config")
		if err != nil {
			if !errors.Is(err, leveldb.ErrNotFound) {
				e(err)
			}
		} else {
			start, err = strconv.Atoi(string(b))
			e(err)
		}
	}
	do(c.MaxID, start, db)
}

func do(maxpage, start int, db *database.LevelDB) {
	a := 0
	w := sync.WaitGroup{}
	for i := start; i <= maxpage; i++ {
		w.Add(1)
		a++
		go func() {
			defer w.Done()
			addonID := strconv.Itoa(i)
			err := save(addonID, db, curseapi.AddonInfo, "modinfo-")
			if err != nil {
				log.Println(addonID, err)
				return
			}
			save(addonID, db, curseapi.Addonfiles, "modfiles-")
		}()
		if a >= 15 {
			w.Wait()
			a = 0
			time.Sleep(time.Second * 1)
			db.Put("config", []byte(strconv.Itoa(i)))
			log.Printf("%d/%d", i, maxpage)
		}
	}
	w.Wait()
}

func save(addonID string, db *database.LevelDB, getfunc func(string) ([]byte, error), keyPrefix string) error {
	_, err := db.Get(keyPrefix + addonID)
	if err != nil {
		if !errors.Is(err, leveldb.ErrNotFound) {
			e(err)
		}
	} else {
		return nil
	}

	var b []byte
	err = retry.Do(func() (err error) {
		b, err = getfunc(addonID)
		return err
	}, retryOpts...)

	if err != nil {
		re := err.(retry.Error)
		for _, v := range re.WrappedErrors() {
			var httperr curseapi.ErrHttpCode
			if errors.As(v, &httperr) {
				if httperr.Code == 404 {
					return fmt.Errorf("save: %w", err)
				}
			}
		}
		panic(err)
	}
	e(db.Put(keyPrefix+addonID, b))
	return nil
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

var retryOpts = []retry.Option{
	retry.Attempts(15),
	retry.Delay(time.Second * 2),
	retry.OnRetry(func(n uint, err error) {
		log.Printf("retry %d: %v", n, err)
	}),
	retry.RetryIf(func(e error) bool {
		var httperr curseapi.ErrHttpCode
		if errors.As(e, &httperr) {
			return httperr.Code != 404
		}
		return true
	}),
}
