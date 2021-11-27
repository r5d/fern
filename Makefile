MOD=ricketyspace.net/fern

fern: fmt
	go build

fmt:
	go fmt ${MOD}
