all: build

build:
	[ -d dist ] || mkdir dist
	go build -o dist/castle

install: build
	sudo mv dist/castle /usr/local/bin/castle

clean:
	[ -d dist ] && rm -rf dist
	[ -f "/usr/local/bin/castle" ] && sudo rm /usr/local/bin/castle