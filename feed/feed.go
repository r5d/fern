// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

package feed

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"ricketyspace.net/fern/schema"
	"ricketyspace.net/fern/state"
)

type Feed struct {
	Id            string
	Source        string
	Schema        string
	Last          int
	TitleContains string `json:"title-contains"`
	YDLPath       string
	DumpDir       string
	Entries       []schema.Entry
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

	// Check 'last'
	if feed.Last < 1 {
		return fmt.Errorf("'last' not set or 0 in a feed '%s'", feed.Id)
	}

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

// Processes the feed.
func (feed *Feed) Process(pState *state.ProcessState) {
	// Init FeedResult.
	fr := state.FeedResult{
		FeedId:     feed.Id,
		FeedResult: "",
		Err:        nil,
	}

	// Get raw feed.
	bs, err := feed.get()
	if err != nil {
		fr.Err = err
		fr.FeedResult = "Unable to get feed"
		pState.FeedResultChan <- fr
		return
	}

	// Unmarshal raw feed into Feed.Object
	err = feed.unmarshal(bs)
	if err != nil {
		fr.Err = err
		fr.FeedResult = "Unable to parse feed"
		pState.FeedResultChan <- fr
		return
	}

	//
	// Process entries.
	//
	// Number entries being processed.
	errors := 0
	processing := 0
	traversed := 0
	// Channel for receiving entry results.
	erChan := make(chan state.EntryResult)
	for _, entry := range feed.Entries {
		e := entry

		// Ignore entry if its title does not matches
		// feed's 'title-contains' string.
		if len(feed.TitleContains) > 0 &&
			!e.TitleContains(feed.TitleContains) {
			fmt.Printf("[%s][%s]: Skipping '%s'\n",
				feed.Id, e.Id, e.Title)
			continue
		}

		// Process entry only if it was not downloaded before.
		if !pState.DB.Exists(feed.Id, e.Id) {
			go feed.processEntry(e, erChan)
			processing += 1
		} else {
			fmt.Printf("[%s][%s]: Already downloaded '%s' before\n",
				feed.Id, e.Id, e.Title)
		}
		traversed += 1

		// Process only `feed.Last` entries.
		if traversed >= feed.Last-1 {
			break
		}
	}
	// Wait for all entries to finish processing.
	for processing > 0 {
		eTxt := "entries"
		if processing == 1 {
			eTxt = "entry"
		}
		fmt.Printf("[%s]: Waiting for %d %s to finish processing\n",
			feed.Id, processing, eTxt)
		er := <-erChan
		if er.Err == nil {
			fmt.Printf("[%s][%s]: Downloaded '%s'\n",
				feed.Id, er.EntryId, er.EntryTitle)
			// Log entry in db.
			pState.DB.Add(feed.Id, er.EntryId)
		} else {
			fmt.Printf("[%s][%s]: Failed to download '%s': %v\n",
				feed.Id, er.EntryId, er.EntryTitle,
				er.Err.Error())
			errors += 1
		}
		processing -= 1
	}
	if errors == 0 {
		fr.FeedResult = "Processed feed"
	} else {
		fr.FeedResult = "Processed feed. One or more" +
			" entries failed to download"
	}
	pState.FeedResultChan <- fr
}

func (feed *Feed) processEntry(entry schema.Entry, erc chan state.EntryResult) {
	// Init EntryResult.
	er := state.EntryResult{
		EntryId:    entry.Id,
		EntryTitle: entry.Title,
		Err:        nil,
	}

	// Download entry.
	fmt.Printf("[%s][%s] Going to download %s\n", feed.Id,
		entry.Id, entry.Title)
	err := feed.ydl(entry.Link)
	if err != nil {
		er.Err = err
	}
	erc <- er
}

func (feed *Feed) ydl(url string) error {
	if len(url) == 0 {
		return fmt.Errorf("URL invalid")
	}

	// Download url via youtube-dl
	outputTemplate := fmt.Sprintf("-o%s",
		path.Join(feed.DumpDir, "%(title)s-%(id)s.%(ext)s"))
	cmd := exec.Command(feed.YDLPath, "--no-progress", outputTemplate, url)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal raw feed into an object.
func (feed *Feed) unmarshal(bs []byte) error {
	var err error

	// Unmarshal based on feed's schema type.
	switch {
	case feed.Schema == "npr":
		feed.Entries, err = nprUnmarshal(bs)
		if err != nil {
			return err
		}
		return nil
	case feed.Schema == "youtube":
		feed.Entries, err = youtubeUnmarshal(bs)
		if err != nil {
			return err
		}
		return nil
	case feed.Schema == "podcast":
		feed.Entries, err = podcastUnmarshal(bs)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("schema of feed '%s' unknown", feed.Id)
}

// Unmarshal a NPR feed.
func nprUnmarshal(bs []byte) ([]schema.Entry, error) {
	nprFeed := new(schema.NPRFeed)
	err := xml.Unmarshal(bs, nprFeed)
	if err != nil {
		return nil, err
	}

	// Get all entries.
	entries := make([]schema.Entry, 0)
	for _, e := range nprFeed.Entries {
		t, err := time.Parse(time.RFC1123Z, e.Pub)
		if err != nil {
			return nil, err
		}
		entry := schema.Entry{
			Id:      e.Id,
			Title:   e.Title,
			PubTime: t,
			Link:    e.Link.Url,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Unmarshal a YouTube feed.
func youtubeUnmarshal(bs []byte) ([]schema.Entry, error) {
	ytFeed := new(schema.YouTubeFeed)
	err := xml.Unmarshal(bs, ytFeed)
	if err != nil {
		return nil, err
	}

	// Get all entries.
	entries := make([]schema.Entry, 0)
	for _, e := range ytFeed.Entries {
		t, err := time.Parse(time.RFC3339, e.Pub)
		if err != nil {
			return nil, err
		}
		entry := schema.Entry{
			Id:      e.Id,
			Title:   e.Title,
			PubTime: t,
			Link:    e.Link.Url,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// Unmarshal a Podcast feed.
func podcastUnmarshal(bs []byte) ([]schema.Entry, error) {
	pcFeed := new(schema.PodcastFeed)
	err := xml.Unmarshal(bs, pcFeed)
	if err != nil {
		return nil, err
	}

	// Get all entries.
	entries := make([]schema.Entry, 0)
	for _, e := range pcFeed.Entries {
		layout := time.RFC1123Z
		if e.Pub[len(e.Pub)-1:] == "T" {
			// Textual time zone. like 'EDT'.
			if e.Pub[6:7] == " " {
				layout = "Mon, 2 Jan 2006 15:04:05 MST"
			} else {
				layout = time.RFC1123
			}
		}
		t, err := time.Parse(layout, e.Pub)
		if err != nil {
			return nil, err
		}
		entry := schema.Entry{
			Id:      e.Id,
			Title:   e.Title,
			PubTime: t,
			Link:    e.Link.Url,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
