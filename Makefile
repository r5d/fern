# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/fern

fern: fmt fix
	go build

fmt:
	go fmt ./...

fix:
	go fix ./...

test:
	go test -v ${MOD}/db
.PHONY: test
