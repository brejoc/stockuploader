package tomlConfig

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type SUConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Remote   string `toml:"remote"`
	Port     string `toml:"port"`
}

func Load(path string) (SUConfig, error) {
	if path == "" {
		return SUConfig{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return SUConfig{}, fmt.Errorf("read config file %q: %w", path, err)
	}
	var c SUConfig
	if err := toml.Unmarshal(data, &c); err != nil {
		return SUConfig{}, fmt.Errorf("parse config %q: %w", path, err)
	}
	return c, nil
}
