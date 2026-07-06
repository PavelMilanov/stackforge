package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/PavelMilanov/stackforge/storage/metadata"
)

/*
ErrTemplateNotFound используется, когда template не найден в прикладном слое.
*/
var ErrTemplateNotFound = errors.New("template not found")

/*
MetadataStore описывает хранилище metadata для custom templates.
*/
type MetadataStore interface {
	GetByTemplateKey(key string) (metadata.TemplateMetadata, error)
}

/*
toServiceInfo преобразует metadata.ServiceInfo в services.ServiceInfo.

Входные параметры:
- items: список сервисов из metadata storage.

Возвращает:
- []ServiceInfo: список сервисов в формате прикладного слоя.
*/
func toServiceInfo(items []metadata.ServiceInfo) []ServiceInfo {
	services := make([]ServiceInfo, 0, len(items))

	for _, item := range items {
		services = append(services, ServiceInfo{
			Name: item.Name,
			Note: item.Note,
		})
	}

	return services
}

/*
TemplatesList возвращает список Portainer templates, дополненный локальной metadata.

Входные параметры:
  - context.Context: контекст вызова; сейчас не используется, так как Portainer client
    работает без внешнего context.

Возвращает:
- []StackTemplate: список templates в формате прикладного слоя.
- error: ошибка Portainer API, ошибка metadata storage или nil при успешной сборке списка.
*/
func (p *PortainerService) TemplatesList(_ context.Context) ([]StackTemplate, error) {
	var stackTemplates []StackTemplate
	items, err := p.client.TemplatesList()
	if err != nil {
		return stackTemplates, err
	}

	for _, item := range items {
		meta, err := p.metadata.GetByTemplateKey(item.Title)
		if err != nil {
			return stackTemplates, err
		}
		stackTemplates = append(stackTemplates, StackTemplate{
			ID:          strconv.Itoa(item.ID),
			Title:       meta.Title,
			Category:    meta.Category,
			Description: meta.Description,
			Fit:         meta.Fit,
			Parameters:  meta.Parameters,
			Services:    toServiceInfo(meta.Services),
		})
	}
	return stackTemplates, nil
}

/*
TemplateGetByID ищет Portainer template по ID и дополняет его локальной metadata.

Входные параметры:
  - context.Context: контекст вызова; сейчас не используется, так как Portainer client
    работает без внешнего context.
  - id: строковый идентификатор Portainer custom template.

Возвращает:
- StackTemplate: найденный template в формате прикладного слоя.
- error: ошибка Portainer API, ошибка metadata storage или nil при успешном поиске.
*/
func (p *PortainerService) TemplateGetByID(_ context.Context, id string) (StackTemplate, error) {
	items, err := p.client.TemplatesList()
	if err != nil {
		return StackTemplate{}, err
	}
	for _, item := range items {
		if strconv.Itoa(item.ID) == id {
			meta, err := p.metadata.GetByTemplateKey(item.Title)
			if err != nil {
				return StackTemplate{}, err
			}
			return StackTemplate{
				ID:          strconv.Itoa(item.ID),
				Title:       meta.Title,
				Category:    meta.Category,
				Description: meta.Description,
				Fit:         meta.Fit,
				Parameters:  meta.Parameters,
				Services:    toServiceInfo(meta.Services),
			}, nil
		}
	}
	return StackTemplate{}, nil
}

/*
StacksList возвращает список стеков Portainer в формате прикладного слоя.

Входные параметры:
  - context.Context: контекст вызова; сейчас не используется, так как Portainer client
    работает без внешнего context.

Возвращает:
- []StackInfo: список стеков с ID, именем и временем создания.
- error: ошибка Portainer API или nil при успешном получении списка.
*/
func (p *PortainerService) StacksList(_ context.Context) ([]StackInfo, error) {
	items, err := p.client.GetStacks()
	if err != nil {
		return nil, err
	}

	stacks := make([]StackInfo, 0, len(items))
	for _, item := range items {
		stacks = append(stacks, StackInfo{
			ID:        item.ID,
			Name:      item.StackName,
			CreatedAt: time.Unix(item.CreationDate, 0),
		})
	}

	return stacks, nil
}
