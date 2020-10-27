UBINDIR ?= /usr/bin
SHAREDIR ?= /usr/share
DESTDIR ?=
GOOS?=linux
GOARCH?=amd64

all: build install

build:
	# go-sqlite require CGO
	CGO_ENABLED=1 go build -o luet-migrate-entropy -tags 'netgo osusergo sqlite_omit_load_extension' -ldflags '-extldflags "-static"' main.go

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 luet-migrate-entropy $(DESTDIR)/$(UBINDIR)/
