// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package state

import "ricketyspace.net/fern/db"

// Contains the result of processing a Feed.
type FeedResult struct {
	FeedId     string // Feed's identifier
	FeedResult string // Feed result
	Err        error  // Set on error
}

// Contains the result of processing an Entry.
type EntryResult struct {
	EntryId    string // Entry's identifier
	EntryTitle string // Entry's title
	Err        error  // Set on error
}

type ProcessState struct {
	DB      *db.FernDB
	// Channel for Feed.Process goroutines to communicate to the
	// caller about the number of entries that are being
	// downloaded for a feed.
	FeedResultChan chan FeedResult
}

func NewProcessState() *ProcessState {
	ps := new(ProcessState)
	ps.FeedResultChan = make(chan FeedResult)
	return ps
}
