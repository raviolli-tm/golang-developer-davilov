// Package config loads YAML configuration files with support for environment
// variable expansion. Values may reference variables using ${NAME} or
// ${NAME:-default} syntax, which allows the same config file to be used both
// locally and inside docker-compose.
package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var envPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-([^}]*))?\}`)

// Load reads the file at path, expands environment variables and unmarshals
// the result into out.
func Load(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}

	expanded := envPattern.ReplaceAllStringFunc(string(data), func(match string) string {
		groups := envPattern.FindStringSubmatch(match)
		if value, ok := os.LookupEnv(groups[1]); ok {
			return value
		}
		return groups[2]
	})

	if err := yaml.Unmarshal([]byte(expanded), out); err != nil {
		return fmt.Errorf("parse config %q: %w", path, err)
	}

	return nil
}
