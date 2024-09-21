# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/fern

fern: fmt fix vet
	go build ${BUILD_OPTS}

fmt:
	go fmt ./...

fix:
	go fix ./...

vet:
	go vet ./...

test:
	go test ${TEST_OPTS} ${MOD}/db ${MOD}/feed ${MOD}/file ${MOD}/schema
.PHONY: test

clean:
	go clean
	rm -f fern-*
	rm -f Makefile.dist*
.PHONY: clean
