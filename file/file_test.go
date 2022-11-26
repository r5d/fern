// SPDX-License-Identifier: ISC
// Copyright © 2022 siddharth <s@ricketyspace.net>

package file

import (
	"bytes"
	"testing"
)

func TestReadFile(t *testing.T) {
	expectedBS := []byte("42 is the answer.\n")
	bs, err := ReadFile("testdata/life.txt")
	if err != nil {
		t.Errorf("read file: %v", err)
		return
	}
	if !bytes.Equal(bs, expectedBS) {
		t.Errorf("read content: '%s' != '%s'", bs, expectedBS)
		return
	}
}
