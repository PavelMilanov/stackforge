package services

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/PavelMilanov/stackforge/storage/metadata"
)

/*
ErrTemplateNotFound используется, когда template не найден в прикладном слое.
*/
var ErrTemplateNotFound = errors.New("template not found")

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
- ctx: контекст вызова, который передается в Portainer API.

Возвращает:
- []StackTemplate: список templates в формате прикладного слоя.
- error: ошибка Portainer API, ошибка metadata storage или nil при успешной сборке списка.
*/
func (p *PortainerService) TemplatesList(ctx context.Context) ([]StackTemplate, error) {
	var stackTemplates []StackTemplate
	items, err := p.client.TemplatesList(ctx)
	if err != nil {
		return stackTemplates, err
	}

	for _, item := range items {
		meta, err := p.metadata.GetTemplate(item.Title)
		if err != nil {
			return stackTemplates, err
		}
		stackTemplates = append(stackTemplates, StackTemplate{
			ID:          strconv.Itoa(item.ID),
			Title:       meta.Title,
			Category:    meta.Category,
			Description: meta.Description,
			Repository:  meta.Repository,
			Metadata:    meta.Metadata,
			Services:    toServiceInfo(meta.Services),
		})
	}
	return stackTemplates, nil
}

/*
TemplateGetByID ищет Portainer template по ID и дополняет его локальной metadata.

Входные параметры:
- ctx: контекст вызова, который передается в Portainer API.
- id: строковый идентификатор Portainer custom template.

Возвращает:
- StackTemplate: найденный template в формате прикладного слоя.
- error: ошибка Portainer API, ошибка metadata storage или nil при успешном поиске.
*/
func (p *PortainerService) TemplateGetByID(ctx context.Context, id string) (StackTemplate, error) {
	items, err := p.client.TemplatesList(ctx)
	if err != nil {
		return StackTemplate{}, err
	}
	for _, item := range items {
		if strconv.Itoa(item.ID) == id {
			meta, err := p.metadata.GetTemplate(item.Title)
			if err != nil {
				return StackTemplate{}, err
			}
			return StackTemplate{
				ID:          strconv.Itoa(item.ID),
				Title:       meta.Title,
				Category:    meta.Category,
				Description: meta.Description,
				Repository:  meta.Repository,
				Metadata:    meta.Metadata,
				Services:    toServiceInfo(meta.Services),
			}, nil
		}
	}
	return StackTemplate{}, nil
}

func (p *PortainerService) StacksByStand(ctx context.Context, standNumber string) ([]StackInfo, error) {
	portainerStacks, err := p.client.GetStacks(ctx)
	if err != nil {
		return nil, err
	}
	stacks := make([]StackInfo, 0, len(portainerStacks))
	templates := make(map[string]metadata.TemplateMetadata)

	for _, item := range portainerStacks {
		stand, stackName := parseStackName(item.StackName)
		if stand != standNumber {
			continue
		}

		meta, ok := templates[stackName]
		if !ok {
			meta, err = p.metadata.GetTemplate(stackName)
			if err != nil {
				return nil, err
			}
			templates[stackName] = meta
		}

		stacks = append(stacks, StackInfo{
			ID:         item.ID,
			Name:       item.StackName,
			CreatedAt:  time.Unix(item.CreationDate, 0),
			Repository: meta.Repository,
			Branch:     "stand/" + stand,
			Domain:     "",
		})
	}
	return stacks, nil
}

/*
parseStackName разделяет имя стека Portainer на номер стенда и имя стека.

Входные параметры:
- name: имя стека в формате <stand>-<stack>.

Возвращает:
- string: номер стенда.
- string: имя стека без номера стенда,.
*/
func parseStackName(name string) (string, string) {
	stand, stackName, ok := strings.Cut(name, "-")
	if !ok {
		return "", name
	}
	return stand, stackName
}
