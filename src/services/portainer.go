package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/PavelMilanov/stackforge/storage/metadata"
)

var ErrTemplateNotFound = errors.New("template not found")

type MetadataStore interface {
	GetByTemplateKey(key string) (metadata.TemplateMetadata, error)
}

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
