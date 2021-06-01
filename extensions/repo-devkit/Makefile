UBINDIR ?= /usr/bin
DESTDIR ?=
EXTNAME := $(shell basename $(shell pwd))

# go tool nm ./luet | grep Commit
override LDFLAGS += -X "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/devkit.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
override LDFLAGS += -X "github.com/Luet-lab/extensions/extensions/repo-devkit/pkg/devkit.BuildCommit=$(shell git rev-parse HEAD)"

all: build install

build:
	CGO_ENABLED=0 go build -o luet-repo-devkit -ldflags '$(LDFLAGS)' luet-repo-devkit.go

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 $(EXTNAME) $(DESTDIR)/$(UBINDIR)/
