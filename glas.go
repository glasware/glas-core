package glas

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/glasware/glas-core/internal"
	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/config"
	"github.com/glasware/glas-core/internal/connection"
	"github.com/glasware/glas-core/internal/connection/telnet"
	"github.com/glasware/glas-core/internal/input"
)

var (
	welcome = strings.Join([]string{
		`\ \      / /__| | ___ ___  _ __ ___   ___  | |_ ___    / ___| | __ _ ___`,
		" \\ \\ /\\ / / _ \\ |/ __/ _ \\| '_ ` _ \\ / _ \\ | __/ _ \\  | |  _| |/ _` / __|",
		`  \ V  V /  __/ | (_| (_) | | | | | |  __/ | || (_) | | |_| | | (_| \__ \_`,
		`   \_/\_/ \___|_|\___\___/|_| |_| |_|\___|  \__\___/   \____|_|\__,_|___( )`,
		`                                                                        |/`,
		`                                         _                      _        _`,
		`  __ _ _ __     _____  ___ __   ___ _ __(_)_ __ ___   ___ _ __ | |_ __ _| |`,
		" / _` | '_ \\   / _ \\ \\/ / '_ \\ / _ \\ '__| | '_ ` _ \\ / _ \\ '_ \\| __/ _` | |",
		`| (_| | | | | |  __/>  <| |_) |  __/ |  | | | | | | |  __/ | | | || (_| | |`,
		` \__,_|_| |_|  \___/_/\_\ .__/ \___|_|  |_|_| |_| |_|\___|_| |_|\__\__,_|_|`,
		`                        |_|`,
		` __  __ _   _ ____         _ _            _`,
		`|  \/  | | | |  _ \    ___| (_) ___ _ __ | |_`,
		`| |\/| | | | | | | |  / __| | |/ _ \ '_ \| __|`,
		`| |  | | |_| | |_| | | (__| | |  __/ | | | |_`,
		`|_|  |_|\___/|____/   \___|_|_|\___|_| |_|\__|`,
	}, "\n")
)

type (
	Glas interface {
		Start(context.Context) error
		SendInput(string)
	}

	glas struct {
		in   chan string
		out  io.Writer
		conn *connection.Connection

		errCh chan error

		cfg *config.Config

		inputHandler input.Handler
	}
)

var (
	_ Glas             = new(glas)
	_ internal.Surface = new(glas)
)

func New(in chan string, out io.Writer, cfgPath string, options ...Option) (Glas, error) {
	g := glas{
		in:  in,
		out: out,
	}

	g.inputHandler = input.New(&g)

	var err error
	g.cfg, err = config.Load(cfgPath, options...)
	if err != nil {
		return nil, fmt.Errorf("config.Load -- %w", err)
	}

	return &g, nil
}

func (g *glas) Start(ctx context.Context) error {
	if err := g.WriteLn(welcome); err != nil {
		return err
	}

	if err := g.WriteLn("The command prefix is set to " + g.cfg.Prefix()); err != nil {
		return err
	}

	return g.eval(ctx)
}

func (g *glas) SendInput(str string) {
	g.in <- str
}

func (g *glas) Echo() bool {
	if conn := g.Connection(); conn != nil {
		return conn.Echo() && g.cfg.Echo()
	}

	return g.cfg.Echo()
}

func (g *glas) CommandPrefix() string {
	return g.cfg.Prefix()
}

func (g *glas) Connection() *connection.Connection {
	return g.conn
}

func (g *glas) NewConnection(ctx context.Context, addr string) error {
	var err error
	g.conn, err = connection.New(addr)
	if err != nil {
		return err
	}

	go g.readSocket(ctx)
	return nil
}

func (g *glas) Write(b []byte) (int, error) {
	return g.out.Write(b)
}

func (g *glas) WriteLn(str string) error {
	_, err := g.Write([]byte(str + "\n"))
	return err
}

func (g *glas) WriteF(format string, v ...any) error {
	str := fmt.Sprintf(format, v...)
	_, err := g.Write([]byte(str))
	return err
}

func (g *glas) Aliases() *actions.Aliases {
	return g.cfg.Aliases()
}

func (g *glas) readSocket(ctx context.Context) {
	if conn := g.Connection(); conn != nil {
		for {
			if !conn.Connected() {
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				byt := make([]byte, 10)
				n, err := conn.Read(byt)
				if err != nil && !errors.Is(err, telnet.ErrClosed) {
					g.errCh <- err
					return
				}

				if n > 0 {
					_, err = g.Write(byt)
					if err != nil {
						g.errCh <- err
						return
					}
				}

				if errors.Is(err, telnet.ErrClosed) {
					return
				}
			}
		}
	}
}
