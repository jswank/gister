package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	flags "github.com/jessevdk/go-flags"
)

const Usage = `[OPTION]... FILE...
  gister [OPTION]... DIR
  CMD | gister [OPTION]...

gister creates GitHub gists using the files specified.

Authentication:
The GitHub API requires authentication.  An OAuth token can be provided via the
environment variable $GISTER_OAUTH_TOKEN or using the configuration file.
`

type Config struct {
	Token       string `description:"OAUTH token for usage with the API" ini-name:"oauth_token" env:"GISTER_OAUTH_TOKEN"`
	Public      bool   `short:"p" long:"public" ini-name:"public" description:"Create a public gist" env:"GISTER_PUBLIC"`
	Description string `short:"d" long:"desc" description:"Description of the gist"`
	Name        string `short:"n" long:"name" description:"Filename if using stdin" default:"stdin.txt"`
	Version     bool   `short:"v" long:"version" description:"Display gister version and exit"`
	Help        bool   `short:"h" long:"help" description:"Show this help message"`
}

var (
	config    = Config{}
	filenames []string
	version   = "2.0"
	parser    = flags.NewParser(&config, flags.PassDoubleDash)
)

func init() {
	parser.Usage = Usage
}

func main() {

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

	filenames, err = getFileList(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate file list, %s\n", err)
		os.Exit(1)
	}

	gist := NewGistCreate()
	gist.Description = config.Description
	gist.Public = config.Public

	for _, f := range filenames {
		gf, err := CreateGistFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding gist, %s\n", err)
			os.Exit(1)
		}
		gist.AddFile(gf)
	}

	// TODO: here on should be put in an improved client interface

	ghc, err := NewGithubClient(config.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client error, %s\n", err)
		os.Exit(2)
	}

	body, err := json.MarshalIndent(gist, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON encoding error, %s\n", err)
		os.Exit(1)
	}

	req, err := ghc.NewAPIRequest("POST", "gists", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error, %s\n", err)
		os.Exit(2)
	}

	res, err := ghc.RunRequest(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connection error, %s\n", err)
		os.Exit(2)
	}

	if res.StatusCode != 201 {
		fmt.Fprintf(os.Stderr, "got bad status from github, %d\n", res.StatusCode)
		data, err := ioutil.ReadAll(res.Body)
		if err == nil {
			fmt.Fprintf(os.Stderr, "%s\n", data)
		}
		os.Exit(2)
	}

	gistResponse, err := parseBody(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing json response: %s\n", err)
		os.Exit(2)
	}

	fmt.Printf("%s\n", gistResponse.HtmlURL)

}

func parseBody(resp io.Reader) (g Gist, err error) {
	data, err := ioutil.ReadAll(resp)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected response: %s\n", data)
	}
	return
}

/* Given a list of potential file names, return a list of files. If a directory
*  is specified, it will be traversed.
 *
 * The size of the args array determines behavior:
 *		empty -> return "-" as the single file to indicate stdin
 *		one element -> either a directory or a single file
 *    larger -> assume the input is the list of files
*/
func getFileList(args []string) ([]string, error) {

	if len(args) < 1 {
		// no args => all files in current dir
		return []string{"-"}, nil
	} else if len(args) == 1 {
		// a file or directory was supplied
		fileinfo, err := os.Stat(args[0])
		if err != nil {
			return nil, err
		}
		if fileinfo.IsDir() {
			// if it is a directory, read it, returning files
			entries, err := os.ReadDir(fileinfo.Name())
			if err != nil {
				return nil, err
			}
			var filelist []string
			for _, e := range entries {
				if e.IsDir() { // skip the directories
					continue
				}
				if e.Type()&os.ModeSymlink == os.ModeSymlink { // skip the symlinks
					continue
				}
				filelist = append(filelist, e.Name())
			}
			return filelist, nil
		} else {
			// otherwise just return the file name
			return []string{args[0]}, nil
		}
	} else {
		return args, nil
	}
}
