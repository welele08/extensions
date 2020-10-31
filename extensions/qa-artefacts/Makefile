UBINDIR ?= /usr/bin
SHAREDIR ?= /usr/share
DESTDIR ?=

all: build install

build:
	CGO_ENABLED=0 go build -o luet-qa-repo-fileconflicts main.go

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 luet-qa-artefacts $(DESTDIR)/$(UBINDIR)/
	install -m 0755 luet-qa-repo-fileconflicts $(DESTDIR)/$(UBINDIR)/