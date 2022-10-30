all: build

build:
	[ -d dist ] || mkdir dist
	go build -o dist/castle

windows:
	[ -d dist ] || mkdir dist
	GOOS=windows go build -o dist/castle.exe

install: build
	sudo cp dist/castle /usr/local/bin/castle

clean:
	[ -d dist ] && rm -rf dist
	[ -f "/usr/local/bin/castle" ] && sudo rm /usr/local/bin/castle

test:
	go test -v ./...