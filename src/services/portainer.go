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
	Status      string
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

func (p *PortainerService) List(_ context.Context) ([]StackTemplate, error) {
	items, err := p.client.GetTemplates()
	if err != nil {
		return nil, err
	}

	var stackTemplates []StackTemplate
	for _, item := range items {
		stackTemplates = append(stackTemplates, StackTemplate{
			ID:          strconv.Itoa(item.ID),
			Name:        item.StackName,
			Category:    "",
			Status:      "",
			Description: "",
			Purpose:     "",
			Fit:         "",
			Parameters:  nil,
			Services:    nil,
		})
	}

	return stackTemplates, nil
}

func (p *PortainerService) GetByID(_ context.Context, id string) (StackTemplate, error) {
	// for _, item := range p.client.GetStacks() {
	// 	if item.ID == id {
	// 		return StackTemplate{
	// 			ID:          item.ID,
	// 			Name:        item.StackName,
	// 			Category:    "",
	// 			Status:      "",
	// 			Description: "",
	// 			Purpose:     "",
	// 			Fit:         "",
	// 			Parameters:  nil,
	// 			Services:    nil,
	// 		}, nil
	// 	}
	return StackTemplate{}, nil
}
