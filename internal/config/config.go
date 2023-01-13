package config

import (
	"fmt"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/actions/aliases"
	"github.com/spf13/afero"
)

const defaultCommandPrefix = "g^"

type Config struct {
	prefix  string
	echo    bool
	aliases *actions.Aliases
	afs     afero.Fs
}

func Load(path string, options ...Option) (*Config, error) {
	cfg := &Config{
		prefix:  defaultCommandPrefix,
		echo:    true,
		aliases: new(actions.Aliases),
	}

	cfg.handleOptions(options...)

	if cfg.afs == nil {
		cfg.afs = afero.NewOsFs()
	}

	s, err := cfg.load(path)
	if err != nil {
		return nil, err
	}

	for _, alias := range s.Aliases {
		if alias.Regex {
			return nil, fmt.Errorf("regex aliases not currently supported")
		}

		cfg.aliases.AddAlias(&aliases.Glob{
			Name:     alias.Name,
			Pattern:  alias.Pattern,
			Template: alias.Template,
		})
	}

	return cfg, nil
}

func (c Config) Prefix() string {
	return c.prefix
}

func (c Config) Echo() bool {
	return c.echo
}

func (c Config) Aliases() *actions.Aliases {
	return c.aliases
}
