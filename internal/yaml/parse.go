package yaml

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func Parse() (*Config, error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("config file not found, using default config")
			return &DefaultConfig, nil
		}
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
