UBINDIR ?= /usr/bin
SHAREDIR ?= /usr/share
DESTDIR ?=

all: build install

build:
	@sed luet-geniso  -e 's|export SCRIPT_DIR.*|export SCRIPT_DIR="$$\{SCRIPT_DIR:-$(DESTDIR)/$(SHAREDIR)/luet-extension/geniso}"|g' -i
	@echo "Nothing to build"

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 luet-geniso $(DESTDIR)/$(UBINDIR)/
	install -d $(DESTDIR)/$(SHAREDIR)/luet-extension/geniso
	install -m 0755 iso/09_prepare.sh $(DESTDIR)/$(SHAREDIR)/luet-extension/geniso
	install -m 0755 iso/10_pack_rootfs.sh $(DESTDIR)/$(SHAREDIR)/luet-extension/geniso
	install -m 0755 iso/13_prepare_iso.sh $(DESTDIR)/$(SHAREDIR)/luet-extension/geniso
	install -m 0755 iso/14_generate_iso.sh $(DESTDIR)/$(SHAREDIR)/luet-extension/geniso
