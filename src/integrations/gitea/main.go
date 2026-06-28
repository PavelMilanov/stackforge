package gitea

import "net/http"

type Client struct {
	Url        string
	Key        string
	httpClient *http.Client
}
