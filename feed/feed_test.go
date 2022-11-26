// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

package feed

import (
	"net/url"
	"testing"

	"ricketyspace.net/fern/file"
)

func TestPodcastUnmarshal(t *testing.T) {
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
		entries, err := podcastUnmarshal(bs)
		if err != nil {
			t.Errorf("feed unmarshal: %v", err)
			return
		}
		for _, entry := range entries {
			if len(entry.Id) < 1 {
				t.Errorf("entry id: %v", entry.Id)
				return
			}
			if len(entry.Title) < 1 {
				t.Errorf("entry title: %v", entry.Title)
				return
			}
			if entry.PubTime.Unix() < 994702392 {
				t.Errorf("entry time: %v", entry.PubTime)
				return
			}
			_, err = url.Parse(entry.Link)
			if err != nil {
				t.Errorf("entry link: %s: %v", entry.Link, err)
				return
			}
		}
	}
}
