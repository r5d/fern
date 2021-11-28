MOD=ricketyspace.net/fern

fern: fmt
	go build

fmt:
	go fmt ${MOD} ${MOD}/config ${MOD}/file ${MOD}/schema
