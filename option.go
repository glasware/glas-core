package glas

import "github.com/glasware/glas-core/internal/config"

type Option = config.Option

func OptNoEcho() Option                     { return config.OptNoEcho() }
func OptCommandPrefix(prefix string) Option { return config.OptCommandPrefix(prefix) }
