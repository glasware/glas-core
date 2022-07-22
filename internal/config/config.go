package config

import (
	"fmt"
	"path/filepath"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/actions/aliases"
)

const defaultCommandPrefix = "g^"

type Config struct {
	CommandPrefix string
	Echo          bool
	Aliases       actions.Aliases
}

func Load(path string, options ...Option) (*Config, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs(%s) -- %w", path, err)
	}

	cfg := new(Config)
	cfg.CommandPrefix = defaultCommandPrefix
	cfg.handleOptions(options...)

	s, err := cfg.load(abs)
	if err != nil {
		return nil, err
	}

	for _, alias := range s.Aliases {
		if alias.Regex {
			return nil, fmt.Errorf("regex aliases not currently supported")
		}

		cfg.Aliases.AddAlias(&aliases.Glob{
			Name:     alias.Name,
			Pattern:  alias.Pattern,
			Template: alias.Template,
		})
	}

	return cfg, nil
}
