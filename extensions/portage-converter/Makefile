UBINDIR ?= /usr/bin
DESTDIR ?=
EXTNAME := $(shell basename $(shell pwd))

# go tool nm ./luet | grep Commit
override LDFLAGS += -X "github.com/Luet-lab/extensions/extensions/portage-converter/cmd.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
override LDFLAGS += -X "github.com/Luet-lab/extensions/extensions/portage-converter/cmd.BuildCommit=$(shell git rev-parse HEAD)"

all: build install

build:
	CGO_ENABLED=0 go build -o luet-portage-converter -ldflags '$(LDFLAGS)'

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 $(EXTNAME) $(DESTDIR)/$(UBINDIR)/
