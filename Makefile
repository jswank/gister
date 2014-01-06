PROJECT  := gister
# turns the git-describe output of v0.3-$NCOMMITS-$SHA1 into
# the more deb friendly 0.3.$NCOMMITS
#VERSION  := $(shell git describe --long --match 'v*' | sed 's/v\([0-9]*\)\.\([0-9]*\)-\([0-9]*\).*/\1.\2.\3/')

VERSION := 1.0

DEBUILD_OPTS = -S -sa

prefix   := /usr
bindir   := $(prefix)/bin
sharedir := $(prefix)/share
mandir   := $(sharedir)/man
man1dir  := $(mandir)/man1

all: build gister.1

build:
	go test
	go build

gister: *.go
	go build
	if [ $(shell basename $(PWD)) != gister ]; then mv $(shell basename $(PWD)) gister; fi

clean:
	rm -rf gister build gister.1 gister.xml gister.html *~

distclean: clean
	rm -rf build/*

gister.html: gister.txt
	asciidoc -b xhtml11 -d manpage -f asciidoc.conf -o $@ $<

gister.xml: gister.txt
	asciidoc -b docbook -d manpage -f asciidoc.conf -o $@ $<

gister.1: gister.xml
	xmlto man $<
	cat $@ | sed -e 's/\[FIXME: source\]/gister/' -e 's/\[FIXME: manual\]/User Commands/' >$@.tmp
	mv $@.tmp $@

install: gister gister.1
	install -D -m 0755 gister $(DESTDIR)$(bindir)/gister
	install -D -m 0644 gister.1 $(DESTDIR)$(man1dir)/gister.1

tar.bz2:
	mkdir -p build
	git archive --prefix="$(PROJECT)-$(VERSION)/" HEAD | bzip2 -z9 >build/$(PROJECT)-$(VERSION).orig.tar.bz2
	git archive --prefix="$(PROJECT)-$(VERSION)/" HEAD | tar -xC build

rpm: tar.bz2
	/usr/bin/rpmdev-setuptree
	cp rpm/gister.spec $(HOME)/rpmbuild/SPECS/
	cp build/$(PROJECT)-$(VERSION).orig.tar.bz2 $(HOME)/rpmbuild/SOURCES/$(PROJECT)-$(VERSION).tar.bz2
	rpmbuild -bb ~/rpmbuild/SPECS/gister.spec

deb: builddeb

builddeb: tar.bz2
	echo $(VERSION) >build/$(PROJECT)-$(VERSION)/version.txt
	(cd build/$(PROJECT)-$(VERSION) && dch --newversion $(VERSION)-1 -b "Last Commit: $(shell git log -1 --pretty=format:'(%ai) %H %cn <%ce>')")
	(cd build/$(PROJECT)-$(VERSION) && dch --release  "new upstream")
	(cd build/$(PROJECT)-$(VERSION) && debuild $(DEBUILD_OPTS) -v$(VERSION)-1)
	@echo "Package is at build/$(PROJECT)-$(VERSION)-1_all.deb"

version:

.PHONY: all build install deb builddeb version clean
