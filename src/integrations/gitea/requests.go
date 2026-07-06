package gitea

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

/*
newRequest создает HTTP-запрос к Gitea API с request timeout и bearer token.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.
- method: HTTP-метод запроса.
- path: путь Gitea API без базового realm.

Возвращает:
- *http.Request: подготовленный HTTP-запрос.
- context.CancelFunc: функция отмены дочернего timeout context.
- error: ошибка создания HTTP-запроса или nil при успешной подготовке.
*/
func (c *Client) newRequest(ctx context.Context, method, path string) (*http.Request, context.CancelFunc, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	timeout := c.timeout
	if timeout <= 0 {
		timeout = defaultRequestTimeout
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	req, err := http.NewRequestWithContext(timeoutCtx, method, c.Realm+path, nil)
	if err != nil {
		cancel()
		return nil, nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)

	return req, cancel, nil
}

/*
ListRepositoryBranches получает список веток репозитория owner/admin из Gitea API.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.

Возвращает:
- []Branch: список веток репозитория.
- error: ошибка создания запроса, HTTP-запроса, чтения тела, статуса ответа или декодирования JSON.
*/
func (c *Client) ListRepositoryBranches(ctx context.Context) ([]Branch, error) {
	req, cancel, err := c.newRequest(ctx, http.MethodGet, "v1/repos/"+c.Owner+"/admin/branches")
	if err != nil {
		return nil, err
	}
	defer cancel()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	var branches []Branch

	err = json.Unmarshal(body, &branches)
	if err != nil {
		return nil, err
	}
	return branches, nil
}
