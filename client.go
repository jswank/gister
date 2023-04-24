package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const GITHUB_API_URL = "https://api.github.com/"

type GistClient struct {
	client *github.Client
}

func NewGistClient(token string) *GistClient {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &GistClient{github.NewClient(tc)}

}

func (gc GistClient) Create(gist *github.Gist) (url string, err error) {

	ctx := context.Background()
	g, _, err := gc.client.Gists.Create(ctx, gist)
	if err != nil {
		return "", err
	}
	return *g.HTMLURL, nil

}
