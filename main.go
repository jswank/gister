package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var (
	filenames []string
	token   string
	version = "1.0"
)

var opts struct {
	Description string `short:"d" long:"desc" description:"Set the description of the gist"`
	Insecure    bool   `short:"i" long:"insecure" description:"Insecure mode; skips SSL cert verification"`
	Public      bool   `short:"p" long:"public" description:"Create a public gist; can be overridden via config"`
	Name        string `short:"s" long:"name" description:"Gist filename (if using stdin)" default:"stdin.txt""`
	ShowVersion bool   `short:"v" long:"version" description:"Display gister and exit"`
}

func init() {


	//
	// parse flags
	//

	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.ShowVersion == true {
		fmt.Printf("%s \n", version)
		os.Exit(0)
	}

	//
	// set token: first try 'git config', then an environment variable
	//

	t, _ := exec.Command("git", "config", "--get", "gist.token").Output()

	token = os.Getenv("GIST_OAUTH")
	if token == "" {
		token = string(bytes.TrimSpace(t))
	}

	if token == "" {
		fmt.Fprintf(os.Stderr, "OAuth token not set; either set git property gist.token or GIST_OAUTH environment variable\n")
		os.Exit(1)
	}

	//
	// use git setting gist.private to determine visibility.
	// 

	p, err := exec.Command("git", "config", "--get", "gist.private").Output()

	if string(bytes.TrimSpace(p)) == "false" {
		opts.Public = true
	}

	//
	// set filenames
	//

	if len(args) < 1 {
		filenames = append(filenames, "-")
	} else {
		filenames = args
	}



}

func main() {

	var gistFiles []GistFile

	for _, f := range filenames {
		g, err := fileToGist(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding gist, %s\n", err)
			os.Exit(1)
		}
		gistFiles = append(gistFiles, g)
	}

	gist := NewGistCreate()
	gist.Description = &opts.Description
	gist.Public = opts.Public

	for _, g := range gistFiles {
		gist.Files[g.Filename] = g
	}

	body, err := json.MarshalIndent(gist, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON encoding error, %s\n", err)
		os.Exit(1)
	}

	ghc, err := NewGithubClient(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client error, %s\n", err)
		os.Exit(2)
	}

	req, err := ghc.NewAPIRequest("POST", "gists", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error, %s\n", err)
		os.Exit(2)
	}

	req.Header.Set("Content-type", "application/json")

	res, err := ghc.RunRequest(req, opts.Insecure)
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

func fileToGist(filename string) (gist GistFile, err error) {

	var input io.Reader

	if filename == "-" { //stdin
		input = os.Stdin
		gist.Filename = opts.Name
	} else {
		input, err = os.Open(filename)
		if err != nil {
			return
		}
		gist.Filename = path.Base(filename)
	}

	file_content, err := ioutil.ReadAll(input)
	if err != nil {
		return
	}

	gist.Content = string(file_content)

	return

}

func parseBody(resp io.Reader) (g GistCreate, err error) {
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
