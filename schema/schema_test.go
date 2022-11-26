// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

package schema

import (
	"encoding/xml"
	"net/url"
	"testing"
	"time"

	"ricketyspace.net/fern/file"
)

func TestPodcastFeed(t *testing.T) {
	testFeeds := []string{
		"testdata/pc-atp.xml",
		"testdata/pc-daringfireball.xml",
		"testdata/pc-kara.xml",
		"testdata/pc-scwpod.xml",
	}
	for _, feed := range testFeeds {
		bs, err := file.ReadFile(feed)
		if err != nil {
			t.Errorf("read feed: %v", err)
			return
		}
		pf := new(PodcastFeed)
		err = xml.Unmarshal(bs, pf)
		if err != nil {
			t.Errorf("xml unmarshal: %v", err)
			return
		}
		for _, entry := range pf.Entries {
			if len(entry.Id) < 1 {
				t.Errorf("entry id: %v", entry.Id)
				return
			}
			if len(entry.Title) < 1 {
				t.Errorf("entry title: %v", entry.Title)
				return
			}
			layout := time.RFC1123Z
			if entry.Pub[len(entry.Pub)-1:] == "T" {
				// Textual time zone. like 'EDT'.
				if entry.Pub[6:7] == " " {
					layout = "Mon, 2 Jan 2006 15:04:05 MST"
				} else {
					layout = time.RFC1123
				}
			}
			pt, err := time.Parse(layout, entry.Pub)
			if err != nil {
				t.Errorf("entry pub: %v: '%v'", layout, entry.Pub)
				return
			}
			if pt.Unix() < 994702392 {
				t.Errorf("entry time: %v", pt)
				return
			}
			_, err = url.Parse(entry.Link.Url)
			if err != nil {
				t.Errorf("entry url: %s: %v", entry.Link.Url, err)
				return
			}
		}
	}
}
