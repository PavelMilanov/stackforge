/*
Package gitea содержит клиент для будущей интеграции с Gitea API.

На текущем этапе пакет хранит конфигурацию клиента, а операции с репозиториями
и ветками будут добавлены по мере реализации workflows создания и удаления стендов.
*/
package gitea

import (
	"net/http"
	"time"
)

const defaultRequestTimeout = 5 * time.Second

/*
Client хранит параметры подключения к Gitea API.
*/
type Client struct {
	Realm      string
	Owner      string
	Token      string
	httpClient *http.Client
	timeout    time.Duration
}

/*
Branch описывает ветку репозитория Gitea.
*/
type Branch struct {
	Name string `json:"name"`
}

/*
NewClient создает HTTP-клиент Gitea API.

Входные параметры:
- realm: базовый URL Gitea API.
- owner: владелец репозитория, с которым работает клиент.
- token: API token для заголовка Authorization.

Возвращает:
- *Client: настроенный клиент Gitea с timeout для каждого запроса.
- error: ошибка инициализации клиента; сейчас всегда nil.
*/
func NewClient(realm, owner, token string) (*Client, error) {
	return &Client{
		Realm:      realm,
		Owner:      owner,
		Token:      token,
		httpClient: &http.Client{},
		timeout:    defaultRequestTimeout,
	}, nil
}
