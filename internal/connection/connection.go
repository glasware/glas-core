package connection

import "io"

//go:generate mockgen -build_flags=--mod=mod -destination=./mock/mock_connection.go -package=mock github.com/glasware/glas-core/internal/connection Connection

type (
	Connection interface {
		Connected() bool
		Close() error
		Addr() string
		Echo() bool
		Write(p []byte) (int, error)
		Read(p []byte) (int, error)
	}
)

var _ io.ReadWriteCloser = Connection(nil)
