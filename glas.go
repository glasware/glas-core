package glas

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/glasware/glas-core/internal/config"
	"github.com/glasware/glas-core/internal/telnet"
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

type Glas struct {
	in         chan string
	out        io.Writer
	socket     *telnet.Socket
	readCancel context.CancelFunc
	errCh      chan error

	cfg *config.Config
}

func New(in chan string, out io.Writer, cfgPath string, options ...Option) (*Glas, error) {
	g := Glas{
		in:  in,
		out: out,
	}

	var err error
	g.cfg, err = config.Load(cfgPath, options...)
	if err != nil {
		return nil, fmt.Errorf("config.Load -- %w", err)
	}

	return &g, nil
}

func (g *Glas) Start(ctx context.Context) error {
	if err := g.writeLn(welcome); err != nil {
		return err
	}

	if err := g.writeLn("The command prefix is set to " + g.cfg.CommandPrefix); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g.errCh = make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case str := <-g.in:
				if err := g.handleInput(ctx, cancel, str); err != nil {
					g.errCh <- err
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		break
	case err := <-g.errCh:
		if err != nil {
			return err
		}
	}

	if g.readCancel != nil {
		g.readCancel()
	}
	return nil
}

func (g *Glas) SendInput(str string) {
	g.in <- str
}

func (g *Glas) writeLn(str string) error {
	return g.write([]byte(str + "\n"))
}

func (g *Glas) writeF(format string, v ...interface{}) error {
	str := fmt.Sprintf(format, v...)
	return g.write([]byte(str))
}

func (g *Glas) write(b []byte) error {
	_, err := g.out.Write(b)
	return err
}

func (g *Glas) readSocket(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			byt := make([]byte, 100)
			_, err := g.socket.Read(byt)
			if err != nil {
				g.errCh <- err
				return
			}

			err = g.write(byt)
			if err != nil {
				g.errCh <- err
				return
			}
		}
	}
}

func (g *Glas) writeToTelnet(b []byte) error {
	b = append(b, []byte("\r\n")...)
	_, err := g.socket.Write(b)
	return err
}
