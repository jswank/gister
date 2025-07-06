package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/go-github/github"
)

var (
	version = "2.1.0"
)

func init() {}

func main() {

	args := doConfig()

	filenames, err := getFileList(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate file list, %s\n", err)
		os.Exit(1)
	}

	gist := github.Gist{
		Description: &config.Description,
		Public:      &config.Public,
	}

	gist.Files = make(map[github.GistFilename]github.GistFile)

	for _, f := range filenames {
		gf, err := createGistFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error encoding gist, %s\n", err)
			os.Exit(1)
		}
		gist.Files[github.GistFilename(*gf.Filename)] = gf
	}

	client := NewGistClient(config.Token)
	url, err := client.Create(&gist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client error, %s\n", err)
		os.Exit(1)
	}
	fmt.Println(url)

}

/* Given a list of potential file names, return a list of files. If a directory
 * is specified, it will be traversed.
 *
 * The size of the args array determines behavior:
 *   empty -> return "-" as the single file to indicate stdin
 *   one element -> either a single file or directory was supplied
 *   larger -> assume the input is the list of files
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

// Given a filename, create a github.GistFile
func createGistFile(filename string) (gf github.GistFile, err error) {

	var input io.Reader

	if filename == "-" { //stdin
		input = os.Stdin
		gf.Filename = &config.Name
	} else {
		input, err = os.Open(filename)
		if err != nil {
			return
		}
		filename = path.Base(filename)
		gf.Filename = &filename
	}

	file_content, err := io.ReadAll(input)
	if err != nil {
		return
	}
	str_content := string(file_content)

	gf.Content = &str_content

	return

}
