install:
	go build
	sudo ln -sf $(shell pwd)/lastpass-search /usr/bin/lastpass-search

