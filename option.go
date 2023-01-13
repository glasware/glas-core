package glas

import (
	"github.com/glasware/glas-core/internal/config"
	"github.com/spf13/afero"
)

type Option = config.Option

func OptNoEcho() Option                     { return config.OptNoEcho() }
func OptCommandPrefix(prefix string) Option { return config.OptCommandPrefix(prefix) }
func OptAfs(afs afero.Fs) Option            { return config.OptAfs(afs) }
