package telnet

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

var ErrClosed = errors.New("connection closed remotely")

type Socket struct {
	conn   net.Conn
	reader *bufio.Reader
	output io.Writer

	Debug bool
}

func (s Socket) Echo() bool {
	return true
}

func (s Socket) Read(p []byte) (int, error) {
	n, err := s.reader.Read(p)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			return n, fmt.Errorf("%w - %s", ErrClosed, err.Error())
		}
	}

	return n, nil
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
