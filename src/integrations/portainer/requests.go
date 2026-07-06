package portainer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

/*
newRequest создает HTTP-запрос к Portainer API с request timeout и API key.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.
- method: HTTP-метод запроса.
- path: путь Portainer API без базового realm.

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
	req.Header.Add("X-API-Key", c.Token)

	return req, cancel, nil
}

/*
GetStacks получает список стеков из Portainer API.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.

Возвращает:
- []Stack: список стеков из Portainer.
- error: ошибка создания запроса, HTTP-запроса, чтения тела, статуса ответа или декодирования JSON.
*/
func (c *Client) GetStacks(ctx context.Context) ([]Stack, error) {
	req, cancel, err := c.newRequest(ctx, http.MethodGet, "/api/stacks")
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

	var stacks []Stack

	err = json.Unmarshal(body, &stacks)
	if err != nil {
		return nil, err
	}
	return stacks, nil
}

/*
GetStackFile получает файл конфигурации для указанного стека.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.
- stack: стек Portainer; для запроса используется поле ID.

Возвращает:
- *Stack: стек с заполненным полем StackFile.
- error: ошибка создания запроса, HTTP-запроса, чтения тела, статуса ответа или декодирования JSON.
*/
func (c *Client) GetStackFile(ctx context.Context, stack Stack) (*Stack, error) {
	req, cancel, err := c.newRequest(ctx, http.MethodGet, "/api/stacks/"+strconv.Itoa(stack.ID)+"/file")
	if err != nil {
		return &stack, err
	}
	defer cancel()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &stack, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &stack, err
	}
	if resp.StatusCode != http.StatusOK {
		return &stack, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	err = json.Unmarshal(body, &stack)
	if err != nil {
		return &stack, err
	}
	return &stack, nil
}

/*
TemplatesList получает список custom templates из Portainer.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.

Возвращает:
- []Template: список custom templates из Portainer.
- error: ошибка создания запроса, HTTP-запроса, чтения тела, статуса ответа или декодирования JSON.
*/
func (c *Client) TemplatesList(ctx context.Context) ([]Template, error) {
	req, cancel, err := c.newRequest(ctx, http.MethodGet, "/api/custom_templates")
	if err != nil {
		return nil, err
	}
	defer cancel()

	query := req.URL.Query()
	query.Set("type", "2") // 2 - docker standalone
	req.URL.RawQuery = query.Encode()

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

	var templates []Template
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

/*
GetTemplateFile получает содержимое файла custom template.

Входные параметры:
- ctx: внешний контекст вызова, от которого наследуется отмена запроса.
- template: custom template Portainer; для запроса используется поле ID.

Возвращает:
- string: содержимое файла custom template.
- error: ошибка создания запроса, HTTP-запроса, чтения тела, статуса ответа или декодирования JSON.
*/
func (c *Client) GetTemplateFile(ctx context.Context, template Template) (string, error) {
	req, cancel, err := c.newRequest(ctx, http.MethodGet, "/api/custom_templates/"+strconv.Itoa(template.ID)+"/file")
	if err != nil {
		return "", err
	}
	defer cancel()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	var data struct {
		FileContent string `json:"FileContent"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	return data.FileContent, nil
}
