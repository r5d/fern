# SPDX-License-Identifier: ISC
# Copyright © 2021 siddharth <s@ricketyspace.net>

MOD=ricketyspace.net/fern

fern: fmt
	go build

fmt:
	go fmt ${MOD} ${MOD}/config ${MOD}/db ${MOD}/feed \
		${MOD}/file ${MOD}/schema ${MOD}/state

test:
	go test -v ${MOD}/db
.PHONY: test
