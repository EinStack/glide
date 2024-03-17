package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type sampleConfig struct {
	Name     string            `yaml:"name"`
	APIKey   string            `yaml:"api_key"`
	Messages map[string]string `yaml:"messages"`
	Seeds    []string          `yaml:"seeds"`
	Params   []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"params"`
}

func TestExpander_EnvVarExpanded(t *testing.T) {
	const apiKey = "ABC1234"

	const seed1 = "40"

	const seed2 = "41"

	const answerMarker = "Answer:"

	const topP = "3"

	const budget = "100"

	t.Setenv("OPENAPI_KEY", apiKey)
	t.Setenv("SEED_1", seed1)
	t.Setenv("SEED_2", seed2)
	t.Setenv("ANSWER_MARKER", answerMarker)
	t.Setenv("OPENAI_TOP_P", topP)
	t.Setenv("OPENAI_BUDGET", budget)

	content, err := os.ReadFile(filepath.Clean(filepath.Join(".", "testdata", "expander.env.yaml")))
	require.NoError(t, err)

	expander := Expander{}
	updatedContent := expander.Expand(content)

	var cfg *sampleConfig

	err = yaml.Unmarshal(updatedContent, &cfg)
	require.NoError(t, err)

	assert.Equal(t, apiKey, cfg.APIKey)
	assert.Equal(t, []string{seed1, seed2, "42"}, cfg.Seeds)

	assert.Contains(t, cfg.Messages["human"], "how $$ $ $ does")
	assert.Contains(t, cfg.Messages["human"], fmt.Sprintf("$%v", answerMarker))

	assert.Equal(t, topP, cfg.Params[0].Value)
	assert.Equal(t, fmt.Sprintf("$%v", budget), cfg.Params[1].Value)
}

func TestExpander_FileContentExpanded(t *testing.T) {
	content, err := os.ReadFile(filepath.Clean(filepath.Join(".", "testdata", "expander.file.yaml")))
	require.NoError(t, err)

	expander := Expander{}
	updatedContent := string(expander.Expand(content))

	require.NotContains(t, updatedContent, "${file:")
	require.Contains(t, updatedContent, "sk-fakeapi-token")
}
