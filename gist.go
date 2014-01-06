package main

// GistCreate represents a gists create payload.
type GistCreate struct {
	URL     string `json:"url,omitempty"`
	HtmlURL string `json:"html_url,omitempty"`
	// the Description is not easily encoded to a string via json.Unmarshal, since null != "" 
	Description *string             `json:"description,omitempty"` // optional
	Public      bool                `json:"public"`                // required
	Files       map[string]GistFile `json:"files"`                 // required
}

// create a new GistDataCreate struct
func NewGistCreate() *GistCreate {
	data := new(GistCreate)
	data.Files = make(map[string]GistFile)

	return data
}

type GistFile struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}
