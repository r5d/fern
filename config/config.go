// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"ricketyspace.net/fern/feed"
	"ricketyspace.net/fern/file"
)

// Represents the fern config
type FernConfig struct {
	YDLPath string      `json:"ydl-path"` // Path to the youtube-dl program.
	DumpDir string      `json:"dump-dir"` // Path where media needs to be downloaded to.
	Feeds   []feed.Feed `json:"feeds"`    // Feeds to download.
}

// Tries to reads the fern config at `$HOME/.config/fern/fern.json`
// and unmarshals it into `FernConfig`
//
// Returns point to `FernConfig` on success. On error, the returned
// config is `nil` and the returned error is non-nil.
func Read() (*FernConfig, error) {
	// Construct config file path.
	h, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	c := path.Join(h, ".config", "fern", "fern.json")

	// Open config file.
	f, err := os.Open(c)
	if err != nil {
		return nil, err
	}

	// Read config file.
	bs, err := file.Read(f)
	if err != nil {
		return nil, err
	}

	// Unmarshal config into an object.
	config := new(FernConfig)
	err = json.Unmarshal(bs, config)
	if err != nil {
		return nil, err
	}

	// Validate config.
	err = config.validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Validates the FernConfig.
//
// Returns nil if validation succeeds; error otherwise.
func (config *FernConfig) validate() error {
	// Validate 'ydl-path' in config.
	if len(config.YDLPath) == 0 {
		return fmt.Errorf("'ydl-path' not set in config")
	}
	_, err := os.Stat(config.YDLPath)
	if err != nil {
		return err
	}

	// Validate 'dump-dir' in config.
	if len(config.DumpDir) == 0 {
		return fmt.Errorf("'dump-dir' not set in config")
	}
	// Replace "~" with user's home directory in the dump
	// directory path.
	h, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	config.DumpDir = strings.Replace(config.DumpDir, "~", h, 1)
	// Ensure dump directory exists.
	err = os.MkdirAll(config.DumpDir, 0755)
	if err != nil {
		return err
	}

	// Validate 'feeds' in config.
	if len(config.Feeds) == 0 {
		return fmt.Errorf("'feeds' not set in config")
	}
	for i, feed := range config.Feeds {
		err = feed.Validate(config.DumpDir)
		if err != nil {
			return err
		}
		config.Feeds[i].YDLPath = config.YDLPath
		config.Feeds[i].DumpDir = feed.DumpDir
	}
	return nil

}
