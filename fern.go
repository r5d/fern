// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

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

	// Initialize process  state.
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
