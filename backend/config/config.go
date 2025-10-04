package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	dotenv "github.com/profclems/go-dotenv"
)

type Env int8

const (
	DEV Env = iota
	PROD
)

var envName = map[Env]string{
	DEV:  "development",
	PROD: "production",
}

func (e Env) String() (string, error) {
	if name, ok := envName[e]; ok {
		return name, nil
	}
	return "unknown", fmt.Errorf("invalid Env value: %d", int(e))
}

func ParseEnv(s string) (Env, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "dev", "development":
		return DEV, nil
	case "prod", "production":
		return PROD, nil
	case "":
		return DEV, nil // default to development
	default:
		return DEV, fmt.Errorf("unknown environment: %s", s)
	}
}

type Config struct {
	Env      Env    `env:"GO_ENV"`
	Database string `env:"DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	envStr := os.Getenv("GO_ENV")
	env, err := ParseEnv(envStr)
	if err != nil {
		return nil, fmt.Errorf("invalid GO_ENV: %w", err)
	}

	if err := loadEnv(env); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := populateWithReflection(config); err != nil {
		return nil, fmt.Errorf("failed to populate config: %w", err)
	}

	if err := exportRuntimeEnv(config); err != nil {
		log.Printf("Warning: failed to export .env.runtime: %v", err)
	}

	return config, nil
}

func loadEnv(env Env) error {
	envString, err := env.String()
	if err != nil {
		panic(fmt.Sprintf("invalid Env value: %d", int(env)))
	}

	// env files to load, latter overrides former
	filesToLoad := []string{
		"../.env",
		fmt.Sprintf("../.env.%s", envString),
		"../.env.local",
	}

	for _, filename := range filesToLoad {
		// Check if file exists before attempting to load
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("Skipping %s (file not found)", filename)
			continue
		}

		// Set the config file and attempt to load
		dotenv.SetConfigFile(filename)
		if err := dotenv.Load(); err != nil {
			return fmt.Errorf("failed to load %s: %w", filename, err)
		}

		log.Printf("Successfully loaded %s", filename)
	}

	return nil
}
