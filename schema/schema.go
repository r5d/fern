// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package schema

import (
	"encoding/xml"
	"time"
)

type Entry struct {
	Id      string
	Title   string
	PubTime time.Time
	Link    string
}

// NPR Feed Schema
type NPRLink struct {
	XMLName xml.Name `xml:"link"`
	Url     string   `xml:",chardata"`
}

type NPREntry struct {
	XMLName xml.Name `xml:"item"`
	Id      string   `xml:"guid"`
	Title   string   `xml:"title"`
	Pub     string   `xml:"pubDate"` // RFC1123Z
	PubTime time.Time
	Link    NPRLink `xml:"link"`
}

type NPRFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Entries []NPREntry `xml:"channel>item"`
}

// YouTube Feed Schema
type YouTubeLink struct {
	XMLName xml.Name `xml:"content"`
	Url     string   `xml:"url,attr"`
}

type YouTubeEntry struct {
	XMLName xml.Name `xml:"entry"`
	Id      string   `xml:"id"`
	Title   string   `xml:"group>title"`
	Pub     string   `xml:"published"` // RFC3339
	PubTime time.Time
	Link    YouTubeLink `xml:"group>content"`
}

type YouTubeFeed struct {
	XMLName xml.Name       `xml:"feed"`
	Entries []YouTubeEntry `xml:"entry"`
}
