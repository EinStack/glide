package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"gopkg.in/yaml.v3"
)

// Provider reads, collects, validates and process config files
type Provider struct {
	expander  *Expander
	Config    *Config
	validator *validator.Validate
}

// NewProvider creates a instance of Config Provider
func NewProvider() *Provider {
	configValidator := validator.New(validator.WithRequiredStructEnabled())

	configValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("yaml"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &Provider{
		expander:  &Expander{},
		Config:    nil,
		validator: configValidator,
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

	err = p.validator.Struct(cfg)

	if err != nil {
		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return p, err
		}

		errors := make([]string, 0, len(err.(validator.ValidationErrors)))

		for _, err := range err.(validator.ValidationErrors) {
			namespace := strings.TrimLeft(err.Namespace(), "Config.")

			errors = append(errors, fmt.Sprintf("- ❌ %v field is %v, \"%v\" provided", namespace, err.Tag(), err.Value()))
		}

		// from here you can create your own error messages in whatever language you wish
		return p, fmt.Errorf(
			"failed to validate config file %v:\n%v\nPlease make sure the config file is properly formatted",
			configPath,
			strings.Join(errors, "\n"),
		)
	}

	p.Config = cfg

	return p, nil
}

func (p *Provider) Get() *Config {
	return p.Config
}

func (p *Provider) Start() {
}
