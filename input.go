package glas

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/glasware/glas-core/internal/telnet"
)

func (g *Glas) handleInput(ctx context.Context, cancel context.CancelFunc, str string) error {
	if g.cfg.Echo {
		if err := g.writeLn(str); err != nil {
			return err
		}
	}

	if strings.HasPrefix(str, g.cfg.CommandPrefix) {
		cmd := strings.TrimPrefix(str, g.cfg.CommandPrefix)
		if err := g.handleCommand(ctx, cancel, cmd); err != nil {
			return err
		}

		return nil
	}

	if g.socket != nil {
		alias, err := g.cfg.Aliases.Check(str)
		if err != nil {
			return err
		}

		if alias != nil {
			var writer io.Writer = g.socket
			if g.cfg.Echo {
				writer = io.MultiWriter(g.socket, g.out)
			}

			return alias(writer)
		}

		return g.writeToTelnet([]byte(str))
	}

	return nil
}

const (
	connect = "connect "
)

func (g *Glas) handleCommand(ctx context.Context, _ context.CancelFunc, cmd string) error {
	switch {
	case strings.HasPrefix(cmd, connect):
		if g.socket != nil {
			if _, err := g.socket.Peek(1); err == nil {
				if err = g.writeF("multiple connections not supported, connected to %s\n", g.socket.Addr()); err != nil {
					return err
				}
				return nil
			}
		}

		addr := strings.TrimPrefix(cmd, connect)

		if err := g.writeLn("connecting to " + addr); err != nil {
			return err
		}

		var err error
		g.socket, err = telnet.Dial(addr)
		if err != nil {
			return fmt.Errorf("telnet.Dial, %s -- %w", addr, err)
		}

		var readCtx context.Context
		readCtx, g.readCancel = context.WithCancel(ctx)
		go g.readSocket(readCtx)

	case cmd == "disconnect":
		if g.readCancel != nil {
			g.readCancel()
		}
	}

	return nil
}
