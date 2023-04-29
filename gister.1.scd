gister(1)

# NAME

gister - Create a githib gist from the command line


# SYNOPSIS

*gister* [OPTION...] FILE...

*gister* [OPTION...] DIR

*gister* -d "gist description" FILE...

*gister* -n "application log" < /var/log/app.log


# DESCRIPTION

gister(1) creates Github gists using one or more files and outputs the
URL of the gist created.  Github OAuth authentication must be used:
see 'AUTHENTICATION' below. 

# OPTIONS

*-d*, *-description*='description'
	Set the description of the gist.

*-h*, *--help*
	Print summary of command line options.

*-n*, *--name*='filename'
	When using stdin, set the file name displayed the gist.  Defaults to "stdin.txt".

*-p*, *--public*
	Create a public gist.  By default, private gists are created.

*-v*, *--version*
	Print gister(1) version.

# AUTHENTICATION

The GitHub API requires authentication. An OAuth token can be provided via the
environment variable $GISTER_OAUTH_TOKEN or using the "oauth_token" value in the
configuration file.

The OAuth token should be authorised to use the "gist" scope with the GithHub
API.

A token can be generataed by visting the URL
https://github.com/settings/tokens/new?scopes=gist

A token can be revoked at anytime by visiting the URL
https://github.com/settings/applications

Once you have an appropriate token, there are two ways to tell gister(1)
what token to use.

*Option 1*: use the configuration directive *oauth_token*

```
$ grep token ~/.config/gister/config
oauth_token = ghp_RRtoaF169Rq0kwmAlZYkkSooxxMEvb123456
```

*Option 2*: specify the token as an environment variable.

```
$ export GISTER_OAUTH_TOKEN="my oauth token"
$ gister example.txt
```

# DEFAULTS

By default, gists created are private.  The *--public* option overrides this
default.

Alternatively, the configuration file can be used to control this behavior:
set *public* to false in order to create public gists by default.

```
$ grep public ~/.config/gister/config
public = false
```

# EXAMPLES

Create a new gist, with description "My test gist", consisting of
two files.

```
$ gister -d "My test gist" example.txt example2.txt
https://gist.github.com/ff9802fdf326abdb6
```

Create a new public gist from standard input.  Assign the name "groups"
to the filename.

```
$ gister -p -n "groups" < /etc/groups
https://gist.github.com/ff9802fdf326abdb6
```

# ENVIRONMENT

*GISTER_OAUTH_CONFIG*
	The token to use for OAuth authentication.  

# FILES
_$XDG_CONFIG_HOME/gister/config_
	An ini-formatted file including configuration directives.

# AUTHOR
gister was originally written by Jason Swank.

# COPYING
Copyright \(C) 2012-2023 Jason Swank.  Free use of this software is
granted under the terms of the MIT license.
