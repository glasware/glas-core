package telnet_test

import (
	"net"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

type server struct {
	listener net.Listener
	conn     net.Conn
}

func newServer(t *testing.T) *server {
	svr := new(server)

	var err error
	svr.listener, err = net.Listen("tcp", "localhost:6666") // FIXME: detect free port.
	require.NoError(t, err)

	go func(tt *testing.T) {
		svr.conn, err = svr.listener.Accept()
		require.NoError(tt, err)
	}(t)

	return svr
}

func (s server) Close() error {
	var merr *multierror.Error

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	if err := s.listener.Close(); err != nil {
		merr = multierror.Append(merr, err)
	}
	return merr.ErrorOrNil()
}
