package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"

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
	rawContent, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return p, fmt.Errorf("unable to read the file %v: %w", configPath, err)
	}

	var rawConfig RawConfig

	if err := yaml.Unmarshal(rawContent, &rawConfig); err != nil {
		return p, fmt.Errorf("unable to serialize the file %v: %w", configPath, err)
	}

	// process raw config
	rawConfig, err = p.expander.Expand(rawConfig)

	if err != nil {
		return p, fmt.Errorf("unable to expand config directives %v: %w", configPath, err)
	}

	// validate the config structure
	var cfg *Config

	err = mapstructure.Decode(rawConfig, cfg)

	if err != nil {
		return p, err
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
