// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

package schema

import (
	"encoding/xml"
	"strings"
	"time"
)

// Generic entry.
type Entry struct {
	Id      string
	Title   string
	PubTime time.Time
	Link    string
}

// Represents a NPR media link.
type NPRLink struct {
	XMLName xml.Name `xml:"link"`
	Url     string   `xml:",chardata"`
}

// Represents an entry in the NPR feed.
type NPREntry struct {
	XMLName xml.Name `xml:"item"`
	Id      string   `xml:"guid"`
	Title   string   `xml:"title"`
	Pub     string   `xml:"pubDate"` // RFC1123Z
	PubTime time.Time
	Link    NPRLink `xml:"link"`
}

// Represents a NPR Feed.
type NPRFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Entries []NPREntry `xml:"channel>item"`
}

// Represents the link a YouTube video.
type YouTubeLink struct {
	XMLName xml.Name `xml:"content"`
	Url     string   `xml:"url,attr"`
}

// Represents an entry in the YouTube feed.
type YouTubeEntry struct {
	XMLName xml.Name `xml:"entry"`
	Id      string   `xml:"id"`
	Title   string   `xml:"group>title"`
	Pub     string   `xml:"published"` // RFC3339
	PubTime time.Time
	Link    YouTubeLink `xml:"group>content"`
}

// Represents a YouTube feed.
type YouTubeFeed struct {
	XMLName xml.Name       `xml:"feed"`
	Entries []YouTubeEntry `xml:"entry"`
}

// Represents a direct link to a Podcast.
type PodcastLink struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

// Represents an entry in the Podcast feed.
type PodcastEntry struct {
	XMLName xml.Name `xml:"item"`
	Id      string   `xml:"guid"`
	Title   string   `xml:"title"`
	Pub     string   `xml:"pubDate"`
	PubTime time.Time
	Link    PodcastLink `xml:"enclosure"`
}

// Represents a iTunes Podcast feed.
type PodcastFeed struct {
	XMLName xml.Name       `xml:"rss"`
	Entries []PodcastEntry `xml:"channel>item"`
}

func (e Entry) TitleContains(contains string) bool {
	return strings.Contains(strings.ToLower(e.Title), strings.ToLower(contains))
}
