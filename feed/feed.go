// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package feed

import (
	"fmt"
	"os"
	"path"
)

type Feed struct {
	Id      string `json:"id"`
	Source  string `json:"source"`
	Schema  string `json:"schema"`
	DumpDir string
}

func (feed *Feed) Validate(baseDumpDir string) error {
	_, err := os.Stat(baseDumpDir)
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

	// Set dump directory for feed and ensure it exists.
	feed.DumpDir = path.Join(baseDumpDir, feed.Id)
	err = os.MkdirAll(feed.DumpDir, 0755)
	if err != nil {
		return err
	}

	return nil
}
