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
BASHCOMPDIR = $(PREFIX)/share/bash-completion/completions
ZSHCOMPDIR = $(PREFIX)/share/zsh/site-functions
FISHCOMPDIR = $(PREFIX)/share/fish/vendor_completions.d

# all: gister completions doc/gister.1
all: gister gister.1

gister:
	$(GO) build $(GOFLAGS)

# completions: gister.bash gister.zsh gister.fish
#
# gister.bash: gister
#	./gister completion bash >gister.bash
#
#gister.zsh: gister
#	./gister completion zsh >gister.zsh
#
#gister.fish: gister
#	./gister completion fish >gister.fish
#
gister.1: gister.1.scd
	$(SCDOC) <gister.1.scd >gister.1

clean:
#	$(RM) -f gister doc/gister.1 gister.bash gister.zsh gister.fish
	$(RM) -f gister doc/gister.1

install:
	$(INSTALL) -d \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/ \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/ \
		$(DESTDIR)$(BASHCOMPDIR) \
		$(DESTDIR)$(ZSHCOMPDIR) \
		$(DESTDIR)$(FISHCOMPDIR)
	$(INSTALL) -pm 0755 gister $(DESTDIR)$(PREFIX)/$(BINDIR)/
	$(INSTALL) -pm 0644 gister.1 $(DESTDIR)$(PREFIX)/$(MANDIR)/man1/
#	$(INSTALL) -pm 0644 gister.bash $(DESTDIR)$(BASHCOMPDIR)/gister
#	$(INSTALL) -pm 0644 gister.zsh $(DESTDIR)$(ZSHCOMPDIR)/_gister
#	$(INSTALL) -pm 0644 gister.fish $(DESTDIR)$(FISHCOMPDIR)/gister.fish

uninstall:
	$(RM) -f \
		$(DESTDIR)$(PREFIX)/$(BINDIR)/gister \
		$(DESTDIR)$(PREFIX)/$(MANDIR)/man1/gister.1 \
		$(DESTDIR)$(BASHCOMPDIR)/gister \
		$(DESTDIR)$(ZSHCOMPDIR)/_gister \
		$(DESTDIR)$(FISHCOMPDIR)/gister.fish

.PHONY: all gister clean install uninstall completions
