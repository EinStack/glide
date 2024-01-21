package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Provider reads, collects, validates and process config files
type Provider struct {
	expander *Expander
	Config   *Config
}

// NewProvider creates a instance of Config Provider
func NewProvider() *Provider {
	return &Provider{
		expander: &Expander{},
		Config:   nil,
	}
}

func (p *Provider) Load(configPath string) (*Provider, error) {
	content, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return p, fmt.Errorf("unable to read config file %v: %w", configPath, err)
	}

	// process raw config
	content = p.expander.Expand(content)

	// validate the config structure
	cfg := DefaultConfig()

	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return p, fmt.Errorf("unable to parse config file %v: %w", configPath, err)
	}

	// TODO: validate config values

	p.Config = cfg

	return p, nil
}

func (p *Provider) Get() *Config {
	return p.Config
}

func (p *Provider) Start() {
}
