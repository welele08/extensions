UBINDIR ?= /usr/bin
SHAREDIR ?= /usr/share
DESTDIR ?=

all: build install

build:
	CGO_ENABLE=0 go build -o luet-migrate-entropy

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 luet-migrate-entropy $(DESTDIR)/$(UBINDIR)/
