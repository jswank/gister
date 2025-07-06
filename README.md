# gister

Create GitHub Gists.

## Quickstart

```console
$ go install scalene.net/gister@latest

$ gister -h
Usage:
  gister [OPTION...] FILE...
  gister [OPTION...] DIR
  CMD | gister [OPTION...]

gister creates GitHub gists using the files specified.

Examples:
  $ gister -d "Description of the gist" file1 file2 file3
  $ tail -100 logs | gister -d "Some logs" -n log.txt

Authentication:
The GitHub API requires authentication. An API token can be provided via the
environment variable $GITHUB_TOKEN or using the configuration file.

Configuration:
A file $XDG_CONFIG_HOME/gister/config can be used to set values for api_token
and public.

Application Options:
  -p, --public   Create a public gist [$GISTER_PUBLIC]
  -d, --desc=    Description of the gist (default: the default description)
  -n, --name=    Filename if using stdin (default: stdin.txt)
  -v, --version  Display gister version and exit
  -h, --help     Show this help message
```

## System Installation
```console
$ make
$ sudo make install
```

## License
MIT
