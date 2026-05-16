package models

type Options struct {
	Daemons map[string]Daemon `yaml:"daemons"`
}

type Daemon struct {
	Target string `yaml:"target"`
}
