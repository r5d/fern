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

type FernConfig struct {
	YDLPath string      `json:"ydl-path"`
	DumpDir string      `json:"dump-dir"`
	Feeds   []feed.Feed `json:"feeds"`
}

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
	for _, feed := range config.Feeds {
		err = feed.Validate(config.DumpDir)
		if err != nil {
			return err
		}
	}
	return nil

}
