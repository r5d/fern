# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/fern

fern: fmt fix vet
	go build

fmt:
	go fmt ./...

fix:
	go fix ./...

vet:
	go vet ./...

test:
	go test -v ${MOD}/db ${MOD}/file ${MOD}/schema
.PHONY: test
