// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

// fern is a simple media feed downloader.
//
// It depends on yt-dlp to download the media found in the media feeds
// to your computer.
//
// fern currently supports YoutTube and NPR feeds.
//
// Information about what media feeds to download, the location of
// yt-dlp program on your computer, and the directory where the media
// should be downloaded to must be specified in a config file which
// fern expects to be at $HOME/.config/fern/fern.json
//
// fern's config file contains three fields:
//
//     {
//        "ydl-path": "/usr/local/bin/yt-dlp",
//        "dump-dir": "~/media/feeds", // media feed download directory
//        "feeds": [...] // list of media feeds.
//     }
//
// Each item in the media "feeds" must be:
//
//     {
//        "id": "media-feed-id", // unique identifier for the media feed
//        "source": "https://feeds.npr.org/XXXX/rss.xml", // media feed url
//        "schema": "npr", // must be "youtube" or "npr"
//        "last": 5 // The last N items that should be downloaded
//     }
//
// You may download an example config file for fern from
// https://ricketyspace.net/fern/fern.json
//
// fern does not take any arguments, to run it just do:
//
//    $ fern
//
package main

import (
	"fmt"
	"os"

	"ricketyspace.net/fern/config"
	"ricketyspace.net/fern/db"
	"ricketyspace.net/fern/state"
)

var fConf *config.FernConfig
var pState *state.ProcessState

func init() {
	var err error

	// Get fern config.
	fConf, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(1)
	}

	// Initialize process state.
	pState = state.NewProcessState()

	// Open database.
	pState.DB, err = db.Open()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	defer pState.DB.Write() // Write database to disk before returning.

	// Process all feeds.
	processing := 0
	for _, feed := range fConf.Feeds {
		f := feed
		go f.Process(pState)
		processing += 1
	}
	// Wait for all feeds finish processing.
	for processing > 0 {
		fTxt := "feeds"
		if processing == 1 {
			fTxt = "feed"
		}
		fmt.Printf("Waiting for %d %s to finish processing\n",
			processing, fTxt)
		fr := <-pState.FeedResultChan
		if fr.Err == nil {
			fmt.Printf("[%s]: %s\n",
				fr.FeedId, fr.FeedResult)
		} else {
			fmt.Printf("[%s]: %s: %v\n",
				fr.FeedId, fr.FeedResult, fr.Err.Error())
		}
		processing -= 1
	}
}
