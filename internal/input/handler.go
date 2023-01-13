package input

import (
	"context"
	"strings"

	"github.com/glasware/glas-core/internal"
)

type Handler struct {
	surface internal.Surface
}

func New(surface internal.Surface) Handler {
	return Handler{
		surface: surface,
	}
}

func (h Handler) Input(ctx context.Context, str string) error {
	if h.surface.Echo() {
		if err := h.surface.WriteLn(str); err != nil {
			return err
		}
	}

	if strings.HasPrefix(str, h.surface.CommandPrefix()) {
		err := h.handleCommand(ctx, strings.TrimPrefix(str, h.surface.CommandPrefix()))
		if err != nil {
			return err
		}

		return nil
	}

	if conn := h.surface.Connection(); conn != nil {
		if conn.Connected() {
			if handled, err := h.handleAlias(str); err != nil {
				return err
			} else if handled {
				return nil
			} else {
				_, err = conn.Write([]byte(str))
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	return nil
}
