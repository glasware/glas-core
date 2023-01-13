package telnet

import (
	"bufio"
	"net"
)

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
