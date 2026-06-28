package portainer

import (
	"net/http"
	"time"
)

/*
PortainerAPI предоставляет интерфейс для взаимодействия с Portainer.
*/
type Client struct {
	Realm      string
	Token      string
	httpClient *http.Client
	Teams      []int
}

/*
Stack представляет абстракцию при взаимодействии с json в Portainer API.
*/
type Stack struct {
	ID           int    `json:"Id"`
	StackName    string `json:"Name"`
	TemplateName string `json:"Title"`
	Endpoint     int    `json:"EndpointId"`
	StackFile    string `json:"StackFileContent"`
	TemplateFile string `json:"FileContent"`
}

// NewPortainer создает новый экземпляр Portainer.
func NewClient(realm, token string, teams []int) (*Client, error) {
	return &Client{
		Realm:      realm,
		Token:      token,
		Teams:      teams,
		httpClient: &http.Client{Timeout: 15 * time.Second}}, nil
}
