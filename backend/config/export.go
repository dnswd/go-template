package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func exportRuntimeEnv(config *Config) error {
	var lines []string

	v := reflect.ValueOf(config).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// Get the actual value from the environment
		envValue := os.Getenv(envTag)

		lines = append(lines, fmt.Sprintf("%s=%s", envTag, envValue))
	}

	content := strings.Join(lines, "\n") + "\n"

	if err := os.WriteFile("../.env.runtime", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write .env.runtime: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Exported %d environment variables to .env.runtime\n", len(lines))
	return nil
}
