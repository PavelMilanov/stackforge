package services

import (
	"context"
	"strconv"

	"github.com/PavelMilanov/stackforge/integrations/portainer"
)

type StackTemplate struct {
	ID          string
	Name        string
	Category    string
	Description string
	Purpose     string
	Fit         string
	Parameters  []string
	Services    []ServiceInfo
}

type PortainerService struct {
	client *portainer.Client
}

func NewPortainerService(client *portainer.Client) *PortainerService {
	return &PortainerService{client: client}
}

func (p *PortainerService) TemplatesList(_ context.Context) ([]StackTemplate, error) {
	items, err := p.client.TemplatesList()
	if err != nil {
		return nil, err
	}

	var stackTemplates []StackTemplate
	for _, item := range items {
		stackTemplates = append(stackTemplates, StackTemplate{
			ID:          strconv.Itoa(item.ID),
			Name:        item.Title,
			Category:    item.Category,
			Description: item.Description,
			Purpose:     item.Note,
			Fit:         "",
			Parameters:  nil,
			Services:    nil,
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
			return StackTemplate{
				ID:          strconv.Itoa(item.ID),
				Name:        item.Title,
				Category:    item.Category,
				Description: item.Description,
				Purpose:     item.Note,
				Fit:         "",
				Parameters:  nil,
				Services:    nil,
			}, nil
		}
	}
	return StackTemplate{}, nil
}
