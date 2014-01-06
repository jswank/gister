package main

import (
	"crypto/tls"
	"io"
	"net/http"
)

const GITHUB_API_URL = "https://api.github.com/"

type GithubClient struct {
	token string
}

func NewGithubClient(token string) (ghc *GithubClient, err error) {
	ghc = &GithubClient{token: token}
	return
}

func (ghc *GithubClient) NewAPIRequest(method, url string, body io.Reader) (req *http.Request, err error) {

	req, err = http.NewRequest(method, GITHUB_API_URL+url, body)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", "token "+ghc.token)

	return
}

func (ghc *GithubClient) RunRequest(req *http.Request, insecure bool) (resp *http.Response, err error) {

	var tr *http.Transport

	if insecure {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{}
	}

	client := &http.Client{Transport: tr}

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}
