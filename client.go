package main

import (
	"io"
	"net/http"
)

const GITHUB_API_URL = "https://api.github.com/"

type GithubClient struct {
	token  string
	client *http.Client
}

func NewGithubClient(token string) (ghc *GithubClient, err error) {
	ghc = &GithubClient{token: token}
	ghc.client = &http.Client{}
	return
}

func (ghc *GithubClient) NewAPIRequest(method, url string, body io.Reader) (req *http.Request, err error) {

	req, err = http.NewRequest(method, GITHUB_API_URL+url, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "token "+ghc.token)
	req.Header.Set("Content-type", "application/json")

	return
}

func (ghc *GithubClient) RunRequest(req *http.Request) (resp *http.Response, err error) {

	resp, err = ghc.client.Do(req)
	if err != nil {
		return
	}

	return
}
