// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"

	"ricketyspace.net/fern/file"
)

var dbPath string

type FernDB struct {
	mutex *sync.Mutex // For writes to `downloaded`
	// Key: feed-id
	// Value: feed-id's entries that were downloaded
	downloaded map[string][]string
}

func init() {
	dbPath = "" // Reset.

	// Construct default dbPath
	h, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dbPath = path.Join(h, ".config", "fern", "db.json")

}

func Open() (*FernDB, error) {
	if len(dbPath) == 0 {
		return nil, fmt.Errorf("FernDB path not set")
	}

	// Check if db exists.
	_, err := os.Stat(dbPath)
	if err != nil {
		// db does not exist yet; create an empty one.
		db := new(FernDB)
		db.mutex = new(sync.Mutex)
		db.downloaded = make(map[string][]string)
		return db, nil
	}

	// Read db from disk.
	f, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	bs, err := file.Read(f)
	if err != nil {
		return nil, err
	}

	// Unmarshal db into an object.
	db := new(FernDB)
	db.mutex = new(sync.Mutex)
	err = json.Unmarshal(bs, &db.downloaded)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (fdb *FernDB) Exists(feed, entry string) bool {
	if _, ok := fdb.downloaded[feed]; !ok {
		return false
	}
	for _, e := range fdb.downloaded[feed] {
		if e == entry {
			return true
		}
	}
	return false

}

func (fdb *FernDB) Add(feed, entry string) {
	// Check if entry already exist for feed.
	if fdb.Exists(feed, entry) {
		return
	}

	// Add entry.
	fdb.mutex.Lock()
	if _, ok := fdb.downloaded[feed]; !ok {
		fdb.downloaded[feed] = make([]string, 0)
	}
	fdb.downloaded[feed] = append(fdb.downloaded[feed], entry)
	fdb.mutex.Unlock()
}

func (fdb *FernDB) Write() error {
	if len(dbPath) == 0 {
		return fmt.Errorf("FernDB path not set")
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Marshal database into json.
	bs, err := json.Marshal(fdb.downloaded)
	if err != nil {
		return err
	}

	// Write to disk.
	_, err = f.Write(bs)
	if err != nil {
		return err
	}
	return nil
}
