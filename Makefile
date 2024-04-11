all:
	CGO_ENABLED=0 go build

man:
	mmark -man ssh2prom.8.md > ssh2prom.8
