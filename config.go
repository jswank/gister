package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// global config variable
var (
	config = Config{}
)

// usage text
const Usage = `[OPTION]... FILE...
  gister [OPTION]... DIR
  CMD | gister [OPTION]...

gister creates GitHub gists using the files specified.

Authentication:
The GitHub API requires authentication.  An OAuth token can be provided via the
environment variable $GISTER_OAUTH_TOKEN or using the configuration file.
`

// The config object, for use by flags
//
// Note - Token can not be set on the command line, and will not appear in the
// -h listing (short/long are not defined).
//
type Config struct {
	Token       string `description:"OAUTH token for usage with the API" ini-name:"oauth_token" env:"GISTER_OAUTH_TOKEN"`
	Public      bool   `short:"p" long:"public" ini-name:"public" description:"Create a public gist" env:"GISTER_PUBLIC"`
	Description string `short:"d" long:"desc" description:"Description of the gist"`
	Name        string `short:"n" long:"name" description:"Filename if using stdin" default:"stdin.txt"`
	Version     bool   `short:"v" long:"version" description:"Display gister version and exit"`
	Help        bool   `short:"h" long:"help" description:"Show this help message"`
}

// do the config init- or die trying.
// return the unparsed args
func doConfig() []string {

	parser := flags.NewParser(&config, flags.PassDoubleDash)
	parser.Usage = Usage

	args, err := parser.Parse()
	if config.Help {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse flags, %s\n", err)
		os.Exit(1)
	}

	if config.Version == true {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	if config.Token == "" {
		fmt.Fprintln(os.Stderr, "OAuth token was not provided")
		os.Exit(1)
	}

	return args

}
