// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package db

import (
	"fmt"
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

func TestExists(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Write a sample test db to fern-db.json
	testDBJSON := []byte(`{"npr":["william-prince","joy-oladokun","lucy-ducas"]}`)
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

	// Test Exists.
	if db.Exists("mkbhd", "v-raptor") {
		t.Errorf("db.Exists failed: mkbhd does not exist in db")
		return
	}
	if db.Exists("npr", "julien-baker") {
		t.Errorf("db.Exists failed: (%s, %s) does not exist in db",
			"npr", "julien-baker")
		return
	}
	if !db.Exists("npr", "joy-oladokun") {
		t.Errorf("db.Exists failed: (%s, %s) exists in db",
			"npr", "joy-oladokun")
		return
	}
}

func TestAdd(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Write a sample test db to fern-db.json
	testDBJSON := []byte(`{"npr":["william-prince","joy-oladokun"]}`)
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

	// Test `Add` for an existing feed.
	db.Add("npr", "julian-baker")
	if len(db.downloaded["npr"]) != 3 {
		t.Errorf("db.Add failed: expected 3 entries for 'npr'")
		return
	}
	if !db.Exists("npr", "julian-baker") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"julian-baker")
		return
	}
	db.Add("npr", "julian-baker")
	if len(db.downloaded["npr"]) != 3 {
		t.Errorf("db.Add failed: expected 3 entries for 'npr'")
		return
	}

	// Test `Add` for nonexistent feed.
	db.Add("mark-rober", "glitter-bomb")
	if len(db.downloaded["mark-rober"]) != 1 {
		t.Errorf("db.Add failed: expected 1 entry for 'mark-rober'")
		return
	}
	if !db.Exists("mark-rober", "glitter-bomb") {
		t.Errorf("db.Add failed: expected %s in 'mark-rober' feed",
			"glitter-bomb")
		return
	}
	db.Add("mark-rober", "squirrel-maze")
	if len(db.downloaded["mark-rober"]) != 2 {
		t.Errorf("db.Add failed: expected 2 entries for 'mark-rober'")
		return
	}
	if !db.Exists("mark-rober", "squirrel-maze") {
		t.Errorf("db.Add failed: expected %s in 'mark-rober' feed",
			"squirrel-maze")
		return
	}
}

func TestWriteNewDB(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Open the db.
	db, err := Open()
	if err != nil {
		t.Errorf("db.Open failed: %v", err.Error())
		return
	}

	// Populate db with test data and write to db to disk.
	db.Add("npr", "william-prince")
	db.Add("npr", "julian-baker")
	db.Add("mkbhd", "v-raptor")
	db.Write()

	// Read db refreshly from disk and verify the db contents.
	db, err = Open()
	if err != nil {
		t.Errorf("db.Open failed: %v", err.Error())
		return
	}
	if len(db.downloaded["npr"]) != 2 {
		t.Errorf("db.Add failed: expected 2 entries for 'npr'")
		return
	}
	if !db.Exists("npr", "william-prince") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"william-prince")
		return
	}
	if !db.Exists("npr", "julian-baker") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"julian-baker")
		return
	}
	if len(db.downloaded["mkbhd"]) != 1 {
		t.Errorf("db.Add failed: expected 1 entry for 'npr'")
		return
	}
	if !db.Exists("mkbhd", "v-raptor") {
		t.Errorf("db.Add failed: expected %s in 'mkbhd' feed",
			"v-raptor")
		return
	}
}

func TestWriteExistingDB(t *testing.T) {
	// Set custom path for db.
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)

	// Write a sample test db to fern-db.json
	testDBJSON := []byte(`{"npr":["kurt-vile","joy-oladokun"]}`)
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

	// Populate db with test data and write to db to disk.
	db.Add("npr", "william-prince")
	db.Add("npr", "julian-baker")
	db.Add("mkbhd", "v-raptor")
	db.Write()

	// Read db refreshly from disk and verify the db contents.
	db, err = Open()
	if err != nil {
		t.Errorf("db.Open failed: %v", err.Error())
		return
	}
	if len(db.downloaded["npr"]) != 4 {
		t.Errorf("db.Add failed: expected 2 entries for 'npr'")
		return
	}
	if !db.Exists("npr", "kurt-vile") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"kurt-vile")
		return
	}
	if !db.Exists("npr", "joy-oladokun") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"joy-oladokun")
		return
	}
	if !db.Exists("npr", "william-prince") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"william-prince")
		return
	}
	if !db.Exists("npr", "julian-baker") {
		t.Errorf("db.Add failed: expected %s in 'npr' feed",
			"julian-baker")
		return
	}
	if len(db.downloaded["mkbhd"]) != 1 {
		t.Errorf("db.Add failed: expected 1 entry for 'npr'")
		return
	}
	if !db.Exists("mkbhd", "v-raptor") {
		t.Errorf("db.Add failed: expected %s in 'mkbhd' feed",
			"v-raptor")
		return
	}
}

func TestConcurrentWrites(t *testing.T) {
	dbPath = path.Join(os.TempDir(), "fern-db.json")
	defer os.Remove(dbPath)
	defer resetDBPath()

	db, err := Open()
	if err != nil {
		t.Errorf("db open failed: %v", err)
		return
	}

	// Randomly create a some entries.
	numEntries := 1000
	entries := make([]string, 0)
	for i := 0; i < numEntries; i++ {
		entries = append(entries, fmt.Sprintf("entry-%d", i))
	}

	// Go routine for adding entries to the db.
	addEntries := func(db *FernDB, feed string, entries []string, donec chan int) {
		for _, entry := range entries {
			db.Add(feed, entry)
		}
		donec <- 1
	}

	// Concurrently write entries to a feed.
	donec := make(chan int)
	feed := "npr"
	routines := 5
	for i := 0; i < routines; i++ {
		go addEntries(db, feed, entries, donec)
	}
	routinesDone := 0
	for routinesDone != routines {
		<-donec
		routinesDone += 1
	}

	// Check if there are exactly numEntries entries.
	if len(db.downloaded[feed]) != numEntries {
		t.Errorf("downloaded entries != %d: %v",
			numEntries, db.downloaded[feed])
	}
}
