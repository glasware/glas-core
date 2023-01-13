package glas

import (
	"context"
)

func (g *glas) eval(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g.errCh = make(chan error, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case str := <-g.in:
				if err := g.inputHandler.Input(ctx, str); err != nil {
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

	if conn := g.Connection(); conn != nil {
		_ = conn.Close()
	}

	return nil
}
