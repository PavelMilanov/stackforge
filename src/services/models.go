package services

import (
	"github.com/PavelMilanov/stackforge/integrations/portainer"
	"github.com/PavelMilanov/stackforge/storage/metadata"
)

type StackTemplate struct {
	ID          string
	Title       string
	Category    string
	Description string
	Fit         string
	Parameters  []string
	Services    []ServiceInfo
}

type ServiceInfo struct {
	Name string
	Note string
}

type PortainerService struct {
	client   *portainer.Client
	metadata MetadataStore
}

func NewPortainerService(client *portainer.Client) *PortainerService {
	store := metadata.NewStore()
	return &PortainerService{client: client, metadata: store}
}
