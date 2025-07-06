package main

import (
	"fmt"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
)

// global config variable
var (
	config = Config{}
)

// usage text
const Usage = `[OPTION...] FILE...
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
and other flags.`

// The config object, for use by flags
//
// Note - Token can not be set on the command line, and will not appear in the
// -h listing (short/long are not defined).
//
type Config struct {
	Token       string `description:"API token for usage with the API" ini-name:"api_token" env:"GITHUB_TOKEN"`
	Public      bool   `short:"p" long:"public" ini-name:"public" description:"Create a public gist (default: false)" env:"GISTER_PUBLIC"`
	Description string `short:"d" long:"desc" description:"Description of the gist" default:"Posted with gister"`
	Name        string `short:"n" long:"name" description:"Filename if using stdin" default:"stdin.txt"`
	Version     bool   `short:"v" long:"version" description:"Display gister version and exit"`
	Help        bool   `short:"h" long:"help" description:"Show this help message"`
}

func init() {}

// do the config init- or die trying.
// return the unparsed args
func doConfig() []string {

	parser := flags.NewParser(&config, flags.PassDoubleDash)
	parser.Usage = Usage

	config_file := getConfigFilename()
	if config_file != "" {
		iniParser := flags.NewIniParser(parser)
		err := iniParser.ParseFile(config_file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse config file, %s\n", err)
			os.Exit(1)
		}

	}
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
		fmt.Fprintln(os.Stderr, "API token was not provided")
		os.Exit(1)
	}

	return args

}

// if an appropriate config file exists, output the filename
func getConfigFilename() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		xdgConfigHome = filepath.Join(os.Getenv("HOME"), ".config")
	}

	config_file := filepath.Join(xdgConfigHome, "gister", "config")

	_, err := os.Stat(config_file)
	if err != nil {
		return ""
	}

	return config_file
}
