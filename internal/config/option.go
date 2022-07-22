package config

type (
	Option interface {
		// Attempt to restrict who can implement this interface.
		// Strictly speaking it doesn't stop it but it should make people
		// think twice, hopefully.
		glasOption() bool
	}

	baseOption      struct{}
	noEchoOption    struct{ baseOption }
	cmdPrefixOption struct {
		baseOption
		prefix string
	}
)

var _ Option = new(baseOption)

func (o baseOption) glasOption() bool { return true }

func OptNoEcho() Option                     { return noEchoOption{} }
func OptCommandPrefix(prefix string) Option { return cmdPrefixOption{prefix: prefix} }

func (c *Config) handleOptions(options ...Option) {
	for _, option := range options {
		switch opt := option.(type) {
		case noEchoOption:
			c.Echo = false
		case cmdPrefixOption:
			c.CommandPrefix = opt.prefix
		}
	}
}
