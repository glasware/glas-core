package telnet

import (
	"bufio"
	"io"
	"net"
)

type Socket struct {
	conn   net.Conn
	reader *bufio.Reader
	output io.Writer
}

func Dial(addr string, options ...Option) (*Socket, error) {

	// FIXME: verify the address string is valid.

	withTimeout, timeout := handleOptions(options...)

	var (
		conn net.Conn
		err  error
	)

	if withTimeout {
		conn, err = net.DialTimeout("tcp", addr, timeout)
	} else {
		conn, err = net.Dial("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	return &Socket{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (s Socket) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}

func (s Socket) Write(p []byte) (int, error) {
	return s.conn.Write(p)
}

func (s Socket) Close() error {
	return s.conn.Close()
}

func (s Socket) Addr() string {
	return s.conn.RemoteAddr().String()
}

func (s Socket) Peek(n int) ([]byte, error) {
	return s.reader.Peek(n)
}
