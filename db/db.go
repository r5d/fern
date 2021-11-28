// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"ricketyspace.net/fern/file"
)

var dbPath string

type FernDB struct {
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
	err = json.Unmarshal(bs, &db.downloaded)
	if err != nil {
		return nil, err
	}
	return db, nil
}
