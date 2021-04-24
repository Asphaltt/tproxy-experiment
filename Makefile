
.PHONY: build clean

build:
	@if [ ! -d bin ]; then mkdir bin; fi
	go build -o bin/tproxy-client ./cmd/tproxy-client
	go build -o bin/tproxy-server ./cmd/tproxy-server
	go build -o bin/tproxy-proxy ./cmd/tproxy-proxy

clean:
	rm -rf bin
