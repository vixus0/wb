go_files := $(shell git ls-files '*.go')

.PHONY: clean

wb wbd: ${go_files}
	CGO_ENABLED=0 go build -o wb ./cmd/wb
	CGO_ENABLED=0 go build -o wbd ./cmd/wbd

clean:
	rm -f wb wbd
