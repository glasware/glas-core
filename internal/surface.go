package internal

import (
	"context"
	"io"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/connection"
)

//go:generate go install github.com/golang/mock/mockgen@latest
//go:generate mockgen -build_flags=--mod=mod -destination=./mock/mock_surface.go -package=mock github.com/glasware/glas-core/internal Surface

import (
	"errors"
)

var ErrExit = errors.New("exit")

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
