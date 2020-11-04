UBINDIR ?= /usr/bin
SHAREDIR ?= /usr/share
DESTDIR ?=
GOOS?=linux
GOARCH?=amd64

all: build install

build:
	# go-sqlite require CGO
	CGO_ENABLED=0 go build -o luet-parallel-tools-commondeps main.go

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 luet-parallel-tools-commondeps $(DESTDIR)/$(UBINDIR)/

