package config

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// Expander finds special directives like ${env:ENV_VAR} in the config file and fill them with actual values
type Expander struct{}

func (e *Expander) Expand(content []byte) []byte {
	expandedContent := string(content)

	expandedContent = e.expandEnvVarDirectives(expandedContent)
	expandedContent = e.expandFileDirectives(expandedContent)
	expandedContent = e.expandEnvVars(expandedContent)

	return []byte(expandedContent)
}

// expandEnvVars expands $ENVAR
func (e *Expander) expandEnvVars(content string) string {
	return os.Expand(content, func(str string) string {
		// This allows escaping environment variable substitution via $$, e.g.
		// - $FOO will be substituted with env var FOO
		// - $$FOO will be replaced with $FOO
		// - $$$FOO will be replaced with $ + substituted env var FOO
		if str == "$" {
			return "$"
		}

		return os.Getenv(str)
	})
}

// expandEnvVarDirectives expands ${env:ENVAR} directives
func (e *Expander) expandEnvVarDirectives(content string) string {
	dirMatcher := regexp.MustCompile(`\$\{env:(.+?)\}`)

	return dirMatcher.ReplaceAllStringFunc(content, func(match string) string {
		matches := dirMatcher.FindStringSubmatch(match)

		if len(matches) != 2 {
			return match // No replacement if the pattern is not matched
		}

		envVarName := matches[1]
		value, exists := os.LookupEnv(envVarName)

		if !exists {
			log.Printf("could not expand the env var directive: \"%s\" variable is not found", envVarName)

			return ""
		}

		return value
	})
}

// expandFileDirectives expands ${file:/path/to/file} directives
func (e *Expander) expandFileDirectives(content string) string {
	dirMatcher := regexp.MustCompile(`\$\{file:(.+?)\}`)

	return dirMatcher.ReplaceAllStringFunc(content, func(match string) string {
		matches := dirMatcher.FindStringSubmatch(match)

		if len(matches) != 2 {
			return match // No replacement if the pattern is not matched
		}

		filePath := matches[1]
		content, err := os.ReadFile(filepath.Clean(filePath))
		if err != nil {
			log.Printf("could not expand the file directive (${file:%s}): %v", filePath, err)
			return match // Return original match if there's an error
		}

		return string(content)
	})
}
