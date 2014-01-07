gister - publish gists from stdin and files
============================================

A command line application to create gists.

Usage
-----
```bash
$ gister < file.txt
$ gister -d "the gist description" -n "filename.txt "< file.txt
$ echo secret | gister
$ echo public | gister -p
$ gister script.py 
$ gister script.js notes.txt
$ pbaste | gist  # Copy from the clipboard - OSX Only
$ xsel | gist  # Copy from the clipboard - X11 Only
$ gister - 
the quick brown fox jumps over the lazy dog
^D
```

Installation
------------

A <a href="http://golang.org/">Go</a> build environment is required to compile.  Follow the instructions in
the <a href="http://golang.org/doc/install">Getting Started</a> document, and configure your <a
href="http://golang.org/doc/articles/go_command.html?h=GOPATH#tmp_3">$GOPATH</a>.  To compile & install, 
just run:
```bash
$ go get github.com/jswank/gister
```

Alternatively, run:

```bash
$ go get github.com/jessevdk/go-flags
$ make
$ sudo make install
```

An RPM can also be created:
```bash
$ make rpm
```

A binary, `gister` and a manual page are created.  This binary can be
installed on any compatible platform and, for normal operation, the
`git` toolchain does not have to be installed in order for it to work.



Authentication
--------------

`gister` uses GitHub OAuth authentication. A token must be generated once and 
then can be used for all your `gister` needs.  To generate an OAuth token, run:

```bash
$ curl -s -u github_username \   
       -d '{"scopes": ["gist"], "note": "commandline gister"}' \     
       https://api.github.com/authorizations
```

The _token_ return value is your token.  This token can be revoked anytime by 
visting https://github.com/settings/applications

There are two ways to tell gister what token to use.

Set the key _gist.token_:
```bash
$ git config --global --add gist.token "my token"
$ gister example.txt
```

Specify the token as an environment variable:
```bash
$ export GIST_OAUTH="your oauth token"
$ gister example.txt
```

Security
--------

By default, `gister` performs validation on the SSL certificate presented by
https://api.github.com.  If your system does not trust the CA certificate used
to sign the certificate presented, then the CA certificate should be added.

Alternatively, the `-i` option can be specified to disable SSL certificate
verification.


Defaults
--------

You can set whether gists are public by using git-config(1) to control 
the behavior.

* _gist.private_ - boolean (true or false) - Determines whether to make a gist public
by default.  If set, then the `--public` option has no effect.


Manual
------

A manual page is available as `gister.txt`; a *NIX manual page is created using this
as a source file as part of the `Makefile`.

License
-------

This program and the accompanying materials are made available under the
terms of the MIT license.
