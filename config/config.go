// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"github.com/elastic/elastic-agent-libs/config"
	"os"
)

type Config struct {
	Dbs      []*config.C `validate:"required"`
	Registry *Registry
}

type Registry struct {
	Path       string      `config:"path"`
	Permission os.FileMode `config:"permission"`
}

type DbExtendConfig struct {
	Addr    string
	Matches []*MatchItemConfig
}

type MatchItemConfig struct {
	Database string
	Tables   []string
}

type name struct {
}

func NewConfig() *Config {
	return &Config{
		Dbs: []*config.C{},
		Registry: &Registry{
			Path:       "position-%s.yml",
			Permission: os.ModePerm,
		},
	}
}
