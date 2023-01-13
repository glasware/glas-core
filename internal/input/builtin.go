package input

import (
	"context"
	"fmt"
	"strings"
)

const (
	connect    = "connect "
	disconnect = "disconnect"
)

const (
	multipleConnectionsNotSupportedFormat = "multiple connections not supported, connected to %s\n"
)

func (h Handler) handleCommand(ctx context.Context, cmd string) error {
	switch {
	case strings.HasPrefix(cmd, connect):
		if err := h.attemptConnect(ctx, strings.TrimPrefix(cmd, connect)); err != nil {
			return err
		}
		return nil

	case cmd == disconnect:
		if conn := h.surface.Connection(); conn != nil {
			conn.Close()
			return nil
		}

		h.surface.WriteLn("no connection")
	}

	return nil
}

func (h Handler) attemptConnect(ctx context.Context, addr string) error {
	if conn := h.surface.Connection(); conn != nil {
		if conn.Connected() {
			if err := h.surface.WriteF(multipleConnectionsNotSupportedFormat, conn.Addr()); err != nil {
				return err
			}
		}
	}

	if err := h.surface.WriteLn("connecting to " + addr); err != nil {
		return err
	}

	if err := h.surface.NewConnection(ctx, addr); err != nil {
		return fmt.Errorf("NewConnection, %s -- %w", addr, err)
	}

	return nil
}
