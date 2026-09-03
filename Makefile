build:
	go build

test:
	go test -v ./...

release:
	env CGO_ENABLED=0 go build -ldflags '-s -w' -trimpath

clean:
	rm -f sidenote

.PHONY: build test release clean
