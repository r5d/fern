// SPDX-License-Identifier: ISC
// Copyright © 2021 siddharth <s@ricketyspace.net>

package file

import (
	"io"
	"os"
)

func Read(f *os.File) ([]byte, error) {
	bs, chunk := make([]byte, 0), make([]byte, 10)
	for {
		n, err := f.Read(chunk)
		if err != nil && err != io.EOF {
			return bs, err
		}
		bs = append(bs, chunk[0:n]...)

		if err == io.EOF {
			break
		}
	}
	return bs, nil
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer f.Close()

	return Read(f)
}

func Write(f *os.File, content []byte) error {
	n, err := f.Write(content)
	if n != len(content) {
		return err
	}
	return nil
}
