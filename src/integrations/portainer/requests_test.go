package portainer

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

/*
roundTripFunc адаптирует функцию под интерфейс http.RoundTripper для тестов.
*/
type roundTripFunc func(*http.Request) (*http.Response, error)

/*
RoundTrip выполняет тестовый HTTP-запрос через функцию roundTripFunc.

Входные параметры:
- r: HTTP-запрос, сформированный тестируемым Portainer client.

Возвращает:
- *http.Response: тестовый HTTP-ответ.
- error: ошибка, которую должен вернуть тестовый transport.
*/
func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

/*
newTestClient создает Portainer client с подмененным HTTP transport.

Входные параметры:
- rt: функция, которая обрабатывает HTTP-запросы в тесте.

Возвращает:
- *Client: тестовый Portainer client с фиксированными realm и token.
*/
func newTestClient(rt roundTripFunc) *Client {
	return &Client{
		Realm:      "http://portainer.test",
		Token:      "test-token",
		httpClient: &http.Client{Transport: rt},
	}
}

/*
TestTemplatesList проверяет получение списка custom templates из Portainer.

Входные параметры:
- t: объект тестирования Go.

Возвращает:
- ничего; при ошибке тест завершает выполнение через методы testing.T.
*/
func TestTemplatesList(t *testing.T) {
	t.Run("returns templates", func(t *testing.T) {
		client := newTestClient(func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/custom_templates" {
				t.Fatalf("path = %s, want /api/custom_templates", r.URL.Path)
			}
			if got := r.URL.Query().Get("type"); got != "2" {
				t.Fatalf("type query = %q, want %q", got, "2")
			}
			if got := r.Header.Get("X-API-Key"); got != "test-token" {
				t.Fatalf("X-API-Key = %q, want %q", got, "test-token")
			}

			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(`[
				{"Id":1,"Title":"Backend","Description":"Backend template","Note":"Use for APIs"},
				{"Id":2,"Title":"Frontend","Description":"Frontend template","Note":"Use for UI"}
			]`)),
			}, nil
		})

		templates, err := client.TemplatesList()
		if err != nil {
			t.Fatalf("TemplatesList returned error: %v", err)
		}

		if len(templates) != 2 {
			t.Fatalf("len(templates) = %d, want 2", len(templates))
		}
		if templates[0].ID != 1 || templates[0].Title != "Backend" ||
			templates[0].Description != "Backend template" || templates[0].Note != "Use for APIs" {
			t.Fatalf("templates[0] = %+v, want backend template", templates[0])
		}
		if templates[1].ID != 2 || templates[1].Title != "Frontend" {
			t.Fatalf("templates[1] = %+v, want frontend template", templates[1])
		}
	})

	t.Run("returns error on non-200 response", func(t *testing.T) {
		client := newTestClient(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "403 Forbidden",
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(strings.NewReader("forbidden")),
			}, nil
		})

		templates, err := client.TemplatesList()
		if err == nil {
			t.Fatal("TemplatesList error = nil, want error")
		}
		if templates != nil {
			t.Fatalf("templates = %+v, want nil", templates)
		}
		if !strings.Contains(err.Error(), "403 Forbidden") || !strings.Contains(err.Error(), "forbidden") {
			t.Fatalf("error = %q, want status and response body", err.Error())
		}
	})

	t.Run("returns error on invalid json", func(t *testing.T) {
		client := newTestClient(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     "200 OK",
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("not json")),
			}, nil
		})

		templates, err := client.TemplatesList()
		if err == nil {
			t.Fatal("TemplatesList error = nil, want error")
		}
		if templates != nil {
			t.Fatalf("templates = %+v, want nil", templates)
		}
	})
}
