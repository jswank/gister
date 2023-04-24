package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	version = "2.0"
)

func init() {}

func main() {

	args := doConfig()

	filenames, err := getFileList(args)
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
 *		one element -> either a single file or directory was supplied
 *    larger -> assume the input is the list of files
*/
func getFileList(args []string) ([]string, error) {

	if len(args) < 1 {
		// no args => all files in current dir
		return []string{"-"}, nil
	} else if len(args) == 1 {
		if args[0] == "-" {
			return []string{args[0]}, nil
		}
		// a file or directory was supplied
		fileinfo, err := os.Stat(args[0])
		if err != nil {
			return nil, err
		}
		// if it is a directory, read it, returning files
		if fileinfo.IsDir() {
			entries, err := os.ReadDir(fileinfo.Name())
			if err != nil {
				return nil, err
			}
			var filelist []string
			for _, e := range entries {
				if e.IsDir() { // skip the directories
					continue
				}
				if e.Type()&os.ModeSymlink == os.ModeSymlink { // skip symlinks
					continue
				}
				filelist = append(filelist, e.Name())
			}
			return filelist, nil
			// otherwise just return the file name
		} else {
			return []string{args[0]}, nil
		}
		// larger -> assume the input is the list of files
	} else {
		return args, nil
	}
}
