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
var defaultDBPath string

// Contains information about list of media that where already
// download for different feeds.
//
// It's stored on disk as a JSON at `$HOME/.config/fern/db.json
type FernDB struct {
	// For locking concurrent read/write access downloaded.
	mutex *sync.RWMutex
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
	defaultDBPath = path.Join(h, ".config", "fern", "db.json")
	dbPath = defaultDBPath
}

// Reads the fern db from disk and unmarshals it into a FernDB
// instance.
//
// Returns a pointer to FernDB on success; nil otherwise. The second
// return value is non-nil on error.
func Open() (*FernDB, error) {
	if len(dbPath) == 0 {
		return nil, fmt.Errorf("FernDB path not set")
	}

	// Check if db exists.
	_, err := os.Stat(dbPath)
	if err != nil {
		// db does not exist yet; create an empty one.
		db := new(FernDB)
		db.mutex = new(sync.RWMutex)
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
	db.mutex = new(sync.RWMutex)
	err = json.Unmarshal(bs, &db.downloaded)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Checks if entry exists in feed. Assumes the current go routine
// already has the mutex lock. Meant for use by the Exists and Add
// methods.
func (fdb *FernDB) exists(feed, entry string) bool {
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

// Returns true if an `entry` for `feed` exists in the database; false
// otherwise.
func (fdb *FernDB) Exists(feed, entry string) bool {
	// Acquire read lock.
	fdb.mutex.RLock()
	defer fdb.mutex.RUnlock() // Give up lock before returning.

	return fdb.exists(feed, entry)
}

// Adds `feed` <-> `entry` to the database.
//
// Once a `feed` <-> `entry` is added to the database, fern assumes
// that entry was downloaded and will not try downloading the entry
// again.
func (fdb *FernDB) Add(feed, entry string) {
	// Acquire write lock.
	fdb.mutex.Lock()
	defer fdb.mutex.Unlock() // Give up lock before returning.

	// Check if entry already exist for feed.
	if fdb.exists(feed, entry) {
		return
	}

	// Add entry.
	if _, ok := fdb.downloaded[feed]; !ok {
		fdb.downloaded[feed] = make([]string, 0)
	}
	fdb.downloaded[feed] = append(fdb.downloaded[feed], entry)

}

// Writes FernDB to disk in the JSON format.
//
// Returns nil on success; error otherwise
func (fdb *FernDB) Write() error {
	// Acquire write lock.
	fdb.mutex.Lock()
	defer fdb.mutex.Unlock() // Give up lock before returning.

	if len(dbPath) == 0 {
		return fmt.Errorf("FernDB path not set")
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

// Sets DB path to the default path. This function is meant to be used
// by tests.
func resetDBPath() {
	dbPath = defaultDBPath
}
