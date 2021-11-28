// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package db

import (
	"os"
	"path"
	"testing"
)

func stringsContain(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func TestOpenPathNotSet(t *testing.T) {
	// Set custom path for db.
	dbPath = ""
	defer os.Remove(dbPath)

	_, err := Open()
	if err == nil {
		t.Errorf("Error: db.Open did not fail when dbPath is empty\n")
		return
	}
	if err.Error() != "FernDB path not set" {
		t.Errorf("Error: db.Open wrong error message when dbPath is empty\n")
		return
	}
}

func TestOpenNewDB(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Open empty db.
	db, err := Open()
	if err != nil {
		t.Errorf("db.Open failed: %v", err.Error())
		return
	}

	// Verify that 'mutex' is initialized.
	if db.mutex == nil {
		t.Errorf("db.mutex is nil")
		return
	}
	db.mutex.Lock()
	db.mutex.Unlock()

	// Verify that 'downloaded' is initialized
	if db.downloaded == nil {
		t.Errorf("db.downloaded is nil")
		return
	}
}

func TestOpenExistingDB(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Write a sample test db to fern-db.json
	testDBJSON := []byte(`{"mkbhd":["rivian","v-raptor","m1-imac"],"npr":["william-prince","joy-oladokun","lucy-ducas"],"simone":["weightless","ugly-desks","safety-hat"]}`)
	dbFile, err := os.Create(dbPath)
	defer dbFile.Close()
	if err != nil {
		t.Errorf("Unable to create fern-db.json: %v", err.Error())
		return
	}
	n, err := dbFile.Write(testDBJSON)
	if len(testDBJSON) != n {
		t.Errorf("Write to fern-db.json failed: %v", err.Error())
		return
	}

	// Open the db.
	db, err := Open()
	if err != nil {
		t.Errorf("db.Open failed: %v", err.Error())
		return
	}

	// Verify that 'mutex' is initialized.
	if db.mutex == nil {
		t.Errorf("db.mutex is nil")
		return
	}
	db.mutex.Lock()
	db.mutex.Unlock()

	// Validate db.downloaded.
	var entries, expectedEntries []string
	var ok bool
	if len(db.downloaded) != 3 {
		t.Errorf("db.downloaded does not contain 3 feeds")
		return
	}
	// mkbhd
	if entries, ok = db.downloaded["mkbhd"]; !ok {
		t.Errorf("db.downloaded does not contain mkbhd")
		return
	}
	expectedEntries = []string{"rivian", "v-raptor", "m1-imac"}
	for _, entry := range entries {
		if !stringsContain(expectedEntries, entry) {
			t.Errorf("%v does not exist in db.downloaded[mkbhd]", entry)
			return
		}
	}
	// simone
	if entries, ok = db.downloaded["simone"]; !ok {
		t.Errorf("db.downloaded does not contain simone")
		return
	}
	expectedEntries = []string{"weightless", "ugly-desks", "safety-hat"}
	for _, entry := range entries {
		if !stringsContain(expectedEntries, entry) {
			t.Errorf("%v does not exist in db.downloaded[simone]", entry)
			return
		}
	}
	// npr
	if entries, ok = db.downloaded["npr"]; !ok {
		t.Errorf("db.downloaded does not contain npr")
		return
	}
	expectedEntries = []string{"william-prince", "lucy-ducas", "joy-oladokun"}
	for _, entry := range entries {
		if !stringsContain(expectedEntries, entry) {
			t.Errorf("%v does not exist in db.downloaded[npr]", entry)
			return
		}
	}
}
