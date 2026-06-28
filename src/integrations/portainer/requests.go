package portainer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

/*
GetStacks получает список стеков из Portainer API.

Returns

	[]Stack - список стеков.
	error - ошибка запроса или декодирования ответа.
*/
func (c *Client) GetStacks() ([]Stack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Realm+"/api/stacks", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-Key", c.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	defer resp.Body.Close()

	var stacks []Stack

	err = json.Unmarshal(body, &stacks)
	if err != nil {
		return nil, err
	}
	return stacks, nil
}

/*
GetStackFile получает файл конфигурации для указанного стека.

Params

	stack - метаданные стека (используется ID).

Returns

	*Stack - стек с заполненным полем StackFile.
	error - ошибка запроса или декодирования ответа.
*/
func (c *Client) GetStackFile(stack Stack) (*Stack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Realm+"/api/stacks/"+strconv.Itoa(stack.ID)+"/file", nil)
	if err != nil {
		return &stack, err
	}
	req.Header.Add("X-API-Key", c.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &stack, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &stack, err
	}
	if resp.StatusCode != http.StatusOK {
		return &stack, fmt.Errorf("%s: %s", resp.Status, string(body))
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &stack)
	if err != nil {
		return &stack, err
	}
	return &stack, nil
}

/*
TemplatesList получает список custom templates из Portainer CustomTemplateList.

Returns

	[]Template - список template-объектов.
	error - ошибка запроса или декодирования ответа.
*/
func (c *Client) TemplatesList() ([]Template, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Realm+"/api/custom_templates", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("type", "2")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("X-API-Key", c.Token)
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

Params

	stack - template-объект (используется ID).

Returns

	*Stack - template-объект с заполненным полем TemplateFile.
	error - ошибка запроса или декодирования ответа.
*/
func (c *Client) GetTemplateFile(stack Stack) (*Stack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Realm+"/api/custom_templates/"+strconv.Itoa(stack.ID)+"/file", nil)
	if err != nil {
		return &stack, err
	}
	req.Header.Add("X-API-Key", c.Token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &stack, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &stack, err
	}
	if resp.StatusCode != http.StatusOK {
		return &stack, fmt.Errorf("%s: %s", resp.Status, string(body))
	}

	defer resp.Body.Close()
	err = json.Unmarshal(body, &stack)
	if err != nil {
		return &stack, err
	}
	return &stack, nil
}
