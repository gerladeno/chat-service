package config

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/gerladeno/chat-service/internal/validator"
)

func ParseAndValidate(filename string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return Config{}, fmt.Errorf("decoding toml config: %w", err)
	}
	if err := validator.Validator.Struct(config); err != nil {
		return Config{}, fmt.Errorf("validationg config: %w", err)
	}
	return config, nil
}
