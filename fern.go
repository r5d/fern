// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

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
//        "last": 5 // the last N items that should be downloaded
//        "title-contains": "tiny desk" // optional. if specified, downloads entries with title matching the value of this field
//     }
//
// You may download an example config file for fern from
// https://ricketyspace.net/fern/fern.json
//
// Run fern with:
//
//    $ fern -run
//
// To print fern's version, do:
//
//    $ fern -version
package main

import (
	"flag"
	"fmt"
	"os"

	"ricketyspace.net/fern/config"
	"ricketyspace.net/fern/db"
	"ricketyspace.net/fern/state"
)

const version = "0.4.0.dev"

var fConf *config.FernConfig
var pState *state.ProcessState

var vFlag *bool
var rFlag *bool

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

	// Parse args.
	vFlag = flag.Bool("version", false, "Print version")
	rFlag = flag.Bool("run", false, "Run fern")
	flag.Parse()

	if *vFlag {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}
	if !*rFlag {
		printUsage(2)
	}
}

func printUsage(exit int) {
	fmt.Printf("fern [ -run | -version ]\n")
	flag.PrintDefaults()
	os.Exit(exit)
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
