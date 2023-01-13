package internal

import (
	"context"
	"io"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/connection"
)

type (
	Surface interface {
		Echo() bool
		CommandPrefix() string

		Connection() *connection.Connection
		NewConnection(ctx context.Context, addr string) error

		io.Writer
		WriteLn(string) error
		WriteF(string, ...any) error

		Aliases() *actions.Aliases
	}
)
