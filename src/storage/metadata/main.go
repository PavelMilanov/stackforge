package metadata

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type TemplateMetadata struct {
	Title       string        `mapstructure:"title"`
	Category    string        `mapstructure:"category"`
	Description string        `mapstructure:"description"`
	Fit         string        `mapstructure:"fit"`
	Parameters  []string      `mapstructure:"parameters"`
	Services    []ServiceInfo `mapstructure:"services"`
}

type ServiceInfo struct {
	Name string `mapstructure:"name"`
	Note string `mapstructure:"note"`
}

var ErrMetadataNotFound = errors.New("template metadata not found")

type Store struct {
	dir string
}

func NewStore() *Store {
	return &Store{dir: filepath.Join("storage", "metadata", "templates")}
}

func (s *Store) GetByTemplateKey(key string) (TemplateMetadata, error) {
	path := filepath.Join(s.dir, key+".yaml")

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return TemplateMetadata{}, fmt.Errorf("template metadata not found: %w", err)
	}

	var meta TemplateMetadata
	if err := v.Unmarshal(&meta); err != nil {
		return TemplateMetadata{}, err
	}

	return meta, nil
}
