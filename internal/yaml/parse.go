package yaml

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("config file not found, using default config")
			return &DefaultConfig, nil
		}
		return nil, err
	}
	return Parse(data)
}

func Parse(data []byte) (*Config, error) {
	if len(data) == 0 {
		return &DefaultConfig, nil
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
