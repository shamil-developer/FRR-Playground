package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"frr-playground/internal/models"
)

func LoadOptions() (models.Options, error) {

	raw, err := os.ReadFile(
		"configs/playground/options.yaml",
	)

	if err != nil {
		return models.Options{}, err
	}

	var options models.Options

	err = yaml.Unmarshal(
		raw,
		&options,
	)

	if err != nil {
		return models.Options{}, err
	}

	return options, nil
}
