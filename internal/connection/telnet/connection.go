package telnet

import (
	"io"

	"github.com/glasware/glas-core/internal/connection"
)

type Connection struct {
	connected bool
	sock      *Socket
}

var _ connection.Connection = Connection{}
var _ io.ReadWriteCloser = Connection{}

func New(addr string) (Connection, error) {
	socket, err := Dial(addr)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		connected: true,
		sock:      socket,
	}, nil
}

func (c Connection) Connected() bool {
	return c.connected
}

func (c Connection) Close() error {
	if !c.connected {
		return nil
	}

	c.connected = false
	return c.sock.Close()
}

func (c Connection) Addr() string {
	if !c.connected {
		return ""
	}

	return c.sock.Addr()
}

func (c Connection) Echo() bool {
	return true // FIXME: should honor telnet negotiation...
}

func (c Connection) Write(p []byte) (int, error) {
	if !c.connected {
		return 0, nil
	}

	p = append(p, []byte("\r\n")...)

	n, err := c.sock.Write(p)
	if err != nil {
		c.connected = false
		return n, err
	}

	return n, nil
}

func (c Connection) Read(p []byte) (int, error) {
	if !c.connected {
		return 0, nil
	}

	n, err := c.sock.Read(p)
	if err != nil {
		c.connected = false
		return n, err
	}

	return n, nil
}

func (c Connection) Peek(n int) ([]byte, error) {
	if !c.connected {
		return nil, nil
	}

	p, err := c.sock.Peek(n)
	if err != nil {
		c.connected = false
		return p, err
	}

	return p, nil
}
