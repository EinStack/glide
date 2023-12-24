package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Provider reads, collects, validates and process config files
type Provider struct {
	Config *Config
}

// NewProvider creates a instance of Config Provider
func NewProvider() *Provider {
	return &Provider{
		Config: nil,
	}
}

func (p *Provider) Load(configPath string) (*Provider, error) {
	rawConfig, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return p, fmt.Errorf("unable to read the file %v: %w", configPath, err)
	}

	var cfg *Config

	if err := yaml.Unmarshal(rawConfig, &cfg); err != nil {
		return p, fmt.Errorf("unable to serialize the file %v: %w", configPath, err)
	}

	p.Config = cfg

	return p, nil
}

func (p *Provider) Get() *Config {
	return p.Config
}

func (p *Provider) Start() {
}
