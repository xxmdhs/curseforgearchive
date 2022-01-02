package main

import (
	"encoding/json"
	"errors"
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

	for _, v := range c.SelID {
		start := 0
		if !v.Reacquire {
			b, err := db.Get("config-" + strconv.Itoa(v.ID))
			if err != nil {
				if !errors.Is(err, leveldb.ErrNotFound) {
					e(err)
				}
			} else {
				start, err = strconv.Atoi(string(b))
				e(err)
			}
		}
		do(v.ID, v.Page, start, db)
	}
}

func do(id, maxpage, start int, db *database.LevelDB) {
	sid := strconv.Itoa(id)
	for i := start; i <= maxpage; i++ {
		var b []byte
		err := retry.Do(func() (err error) {
			b, err = curseapi.Searchmod("", i, id)
			return err
		}, retryOpts...)
		e(err)

		var list []interface{}
		e(json.Unmarshal(b, &list))
		idList := make([]string, 0, len(list))
		for _, v := range list {
			id := strconv.FormatFloat(v.(map[string]interface{})["id"].(float64), 'f', -1, 64)
			idList = append(idList, id)
			toSave("modinfo-", id, v, db)
		}
		get(idList, db)

		db.Put("config-"+sid, []byte(strconv.Itoa(i)))
		log.Printf("id: %d, page: %d, maxpage: %d", id, i, maxpage)
	}
}

func toSave(key, id string, data interface{}, db *database.LevelDB) {
	b, err := json.Marshal(data)
	e(err)
	err = db.Put(key+id, b)
	e(err)
}

func get(l []string, db *database.LevelDB) error {
	i := 0
	w := sync.WaitGroup{}
	for _, v := range l {
		w.Add(1)
		i++
		go func(v string) {
			defer w.Done()
			var b []byte
			err := retry.Do(func() (err error) {
				b, err = curseapi.Addonfiles(v)
				return err
			}, retryOpts...)
			e(err)
			e(db.Put("modfiles-"+v, b))
		}(v)
		if i >= 15 {
			w.Wait()
			i = 0
			time.Sleep(time.Second * 1)
		}
	}
	w.Wait()
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
}
