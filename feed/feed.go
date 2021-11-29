// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package feed

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"ricketyspace.net/fern/schema"
)

type Feed struct {
	Id      string `json:"id"`
	Source  string `json:"source"`
	Schema  string `json:"schema"`
	YDLPath string
	DumpDir string
	Object  interface{}
}

func (feed *Feed) Validate(ydlPath, baseDumpDir string) error {
	_, err := os.Stat(ydlPath)
	if err != nil {
		return err
	}
	_, err = os.Stat(baseDumpDir)
	if err != nil {
		return err
	}

	// Check 'id'
	if len(feed.Id) == 0 {
		return fmt.Errorf("'id' not set in a feed")
	}

	// Check 'source'
	if len(feed.Source) == 0 {
		return fmt.Errorf("'source' not set in a feed '%s'", feed.Id)
	}

	// Check 'schema'
	schemaOK := false
	for _, schema := range []string{"npr", "youtube"} {
		if feed.Schema == schema {
			schemaOK = true
		}
	}
	if !schemaOK {
		return fmt.Errorf("schema '%s' for feed '%s' is not valid",
			feed.Schema, feed.Id)
	}

	// Set ydl-path for feed.
	feed.YDLPath = ydlPath

	// Set dump directory for feed and ensure it exists.
	feed.DumpDir = path.Join(baseDumpDir, feed.Id)
	err = os.MkdirAll(feed.DumpDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

// Get the feed.
func (feed *Feed) get() ([]byte, error) {
	// Init byte container to store feed content.
	bs := make([]byte, 0)

	resp, err := http.Get(feed.Source)
	if err != nil {
		return bs, err
	}

	// Slurp body.
	chunk := make([]byte, 100)
	for {
		c, err := resp.Body.Read(chunk)
		if c < 1 {
			break
		}
		if err != nil && err != io.EOF {
			return bs, err
		}
		bs = append(bs, chunk[0:c]...)
	}
	return bs, nil
}


// Unmarshal raw feed into an object.
func (feed *Feed) unmarshal(bs []byte) error {
	var err error

	// Unmarshal based on feed's schema type.
	switch {
	case feed.Schema == "npr":
		feed.Object, err = nprUnmarshal(bs)
		if err != nil {
			return err
		}
		return nil
	case feed.Schema == "youtube":
		feed.Object, err = youtubeUnmarshal(bs)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("schema of feed '%s' unknown", feed.Id)
}

// Unmarshal a NPR feed.
func nprUnmarshal(bs []byte) (schema.NPRFeed, error) {
	nprFeed := new(schema.NPRFeed)
	err := xml.Unmarshal(bs, nprFeed)
	if err != nil {
		return *nprFeed, err
	}

	// Parse time for all entries.
	for i, entry := range nprFeed.Entries {
		nprFeed.Entries[i].PubTime, err = time.Parse(time.RFC1123Z, entry.Pub)
		if err != nil {
			return *nprFeed, err
		}
	}
	return *nprFeed, nil
}

// Unmarshal a YouTube feed.
func youtubeUnmarshal(bs []byte) (schema.YouTubeFeed, error) {
	ytFeed := new(schema.YouTubeFeed)
	err := xml.Unmarshal(bs, ytFeed)
	if err != nil {
		return *ytFeed, err
	}

	// Parse time for all entries.
	for i, entry := range ytFeed.Entries {
		ytFeed.Entries[i].PubTime, err = time.Parse(time.RFC3339, entry.Pub)
		if err != nil {
			return *ytFeed, err
		}
	}
	return *ytFeed, nil
}
