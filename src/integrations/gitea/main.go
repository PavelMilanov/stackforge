package gitea

import "net/http"

/*
Client хранит параметры подключения к Gitea API.
*/
type Client struct {
	Url        string
	Key        string
	httpClient *http.Client
}
