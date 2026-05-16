package models

type Plan struct {
	Steps []Step `yaml:"steps"`
}

type Step struct {
	Type        string         `yaml:"type"`
	Daemon      string         `yaml:"daemon"`
	Description string         `yaml:"description"`
	Params      map[string]any `yaml:"params"`
}
