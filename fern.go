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

	fConf, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(1)
	}

	pState = state.NewProcessState()
	pState.YDLPath = fConf.YDLPath
	pState.DumpDir = fConf.DumpDir

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
	for _, feed := range fConf.Feeds {
		f := feed
		go f.Process(pState)
		pState.FeedsProcessing += 1
	}
	// Wait for all feeds finish processing.
	for pState.FeedsProcessing > 0 {
		fr := <-pState.FeedResultChan
		if fr.Err == nil {
			fmt.Printf("[%s]: %s\n",
				fr.FeedId, fr.FeedResult)
		} else {
			fmt.Printf("[%s]: %s: %v\n",
				fr.FeedId, fr.FeedResult, fr.Err.Error())
		}
		pState.FeedsProcessing -= 1
	}
}
