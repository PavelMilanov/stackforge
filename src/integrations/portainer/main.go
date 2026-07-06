/*
Package portainer содержит HTTP-клиент для работы с Portainer API.

Пакет отвечает только за сетевые запросы, авторизацию через X-API-Key,
декодирование JSON-ответов и возврат низкоуровневых моделей Portainer.
*/
package portainer

import (
	"net/http"
	"time"
)

const defaultRequestTimeout = 5 * time.Second

/*
Client хранит параметры подключения к Portainer API.
*/
type Client struct {
	Realm      string
	Token      string
	httpClient *http.Client
	timeout    time.Duration
	Teams      []int
}

/*
Stack описывает стек Portainer и поля JSON, которые используются приложением.
*/
type Stack struct {
	ID           int    `json:"Id"`
	StackName    string `json:"Name"`
	TemplateName string `json:"Title"`
	Endpoint     int    `json:"EndpointId"`
	CreationDate int64  `json:"CreationDate"`
	StackFile    string `json:"StackFileContent"`
	TemplateFile string `json:"FileContent"`
}

/*
Template описывает custom template из Portainer API.
*/
type Template struct {
	ID          int    `json:"Id"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Category    string `json:"Logo"`
	Note        string `json:"Note"`
}

/*
NewClient создает HTTP-клиент Portainer API.

Входные параметры:
- realm: базовый URL Portainer без завершающего API path.
- token: API key для заголовка X-API-Key.
- teams: список team ID, доступный клиенту для будущих операций.

Возвращает:
- *Client: настроенный клиент Portainer с timeout для каждого запроса.
- error: ошибка инициализации клиента; сейчас всегда nil.
*/
func NewClient(realm, token string, teams []int) (*Client, error) {
	return &Client{
		Realm:      realm,
		Token:      token,
		Teams:      teams,
		httpClient: &http.Client{},
		timeout:    defaultRequestTimeout,
	}, nil
}
