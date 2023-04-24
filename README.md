# gister

Create GitHub Gists.

## Quickstart

```console
$ gister -h
Usage:
  gister [OPTION]... FILE...
  gister [OPTION]... DIR
  CMD | gister [OPTION]...

gister creates GitHub gists using the files specified.

Authentication:
The GitHub API requires authentication.  An OAuth token can be provided via the
environment variable $GISTER_OAUTH_TOKEN or using the configuration file.


Application Options:
  -p, --public   Create a public gist [$GISTER_PUBLIC]
  -d, --desc=    Description of the gist
  -n, --name=    Filename if using stdin (default: stdin.txt)
  -v, --version  Display gister version and exit
  -h, --help     Show this help message
```
