package services

import (
	"time"

	"github.com/PavelMilanov/stackforge/integrations/portainer"
	"github.com/PavelMilanov/stackforge/storage/metadata"
)

/*
StackTemplate описывает custom template в формате, который используется
прикладным слоем и страницами приложения.
*/
type StackTemplate struct {
	ID          string
	Title       string
	Category    string
	Description string
	Repository  string
	Metadata    []string
	Services    []ServiceInfo
}

/*
ServiceInfo описывает сервис, входящий в template metadata.
*/
type ServiceInfo struct {
	Name string
	Note string
}

/*
StackInfo описывает стек Portainer в формате, который нужен странице стендов.
*/
type StackInfo struct {
	ID         int
	Name       string
	CreatedAt  time.Time
	Repository string
	Branch     string
	Domain     string
}

/*
PortainerService объединяет Portainer API и metadata storage.
*/
type PortainerService struct {
	client   *portainer.Client
	metadata MetadataStore
}

/*
MetadataStore описывает хранилище metadata для custom templates.
*/
type MetadataStore interface {
	GetTemplate(key string) (metadata.TemplateMetadata, error)
}

/*
NewPortainerService создает сервис для работы с Portainer templates и stacks.

Входные параметры:
- client: HTTP-клиент Portainer API.

Возвращает:
- *PortainerService: сервис с подключенным Portainer client и хранилищем metadata.
*/
func NewPortainerService(client *portainer.Client) *PortainerService {
	store := metadata.NewStore()
	return &PortainerService{client: client, metadata: store}
}
