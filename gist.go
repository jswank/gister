package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// GistCreate represents a gists create payload.
type Gist struct {
	URL     string `json:"url,omitempty"`
	HtmlURL string `json:"html_url,omitempty"`
	// the Description is not easily encoded to a string via json.Unmarshal, since null != ""
	Description string              `json:"description,omitempty"` // optional
	Public      bool                `json:"public"`                // required
	Files       map[string]GistFile `json:"files"`                 // required
}

// create a new GistDataCreate struct
func NewGistCreate() *Gist {
	gist := Gist{}
	gist.Files = make(map[string]GistFile)

	return &gist
}

func (g *Gist) AddFile(gf GistFile) {
	g.Files[gf.Filename] = gf
}

type GistFile struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

func CreateGistFile(filename string) (gist GistFile, err error) {

	var input io.Reader

	if filename == "-" { //stdin
		input = os.Stdin
		gist.Filename = config.Name
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
