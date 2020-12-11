BUNDLED_EXTENSIONS = geninitramfs kernel-switcher autobump-github geniso genimage qa-artefacts migrate-entropy package-browser parallel-tools portage-converter portage apkbuildconverter
UBINDIR ?= /usr/bin
DESTDIR ?=

all: build

build:
	for d in $(BUNDLED_EXTENSIONS); do $(MAKE) -C extensions/$$d build; done

install: build
	for d in $(BUNDLED_EXTENSIONS); do $(MAKE) -C extensions/$$d install; done
