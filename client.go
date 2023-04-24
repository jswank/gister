package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GistClient struct {
	client *github.Client
}

// Create a new GistClient using the supplied OAuth token
func NewGistClient(token string) *GistClient {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &GistClient{github.NewClient(tc)}

}

// Create a new gist, returning the URL
func (gc GistClient) Create(gist *github.Gist) (url string, err error) {

	ctx := context.Background()
	g, _, err := gc.client.Gists.Create(ctx, gist)
	if err != nil {
		return "", err
	}
	return *g.HTMLURL, nil

}
