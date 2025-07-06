.POSIX:
.SUFFIXES:

GO = go
RM = rm
INSTALL = install
SCDOC = scdoc
GOFLAGS =
PREFIX = /usr/local
BINDIR = bin
MANDIR = share/man

all: gister gister.1

gister: *.go go.mod go.sum
	$(GO) build $(GOFLAGS)

gister.1: gister.1.scd
	$(SCDOC) <gister.1.scd >gister.1

clean:
	$(RM) -f gister doc/gister.1

install:
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/ 
	$(INSTALL) -pm 0755 gister $(DESTDIR)$(PREFIX)/$(BINDIR)/
	$(INSTALL) -pm 0644 gister.1 $(DESTDIR)$(PREFIX)/$(MANDIR)/man1/

uninstall:
	$(RM) -f \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/gister \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/gister.1 

.PHONY: all clean install uninstall
