package input

import "io"

func (h Handler) handleAlias(str string) (bool, error) {
	alias, err := h.surface.Aliases().Check(str)
	if err != nil {
		return false, err
	}

	if alias != nil {
		var writer io.Writer = h.surface.Connection()
		if h.surface.Echo() {
			writer = io.MultiWriter(writer, h.surface)
		}
		return true, alias(writer)
	}

	return false, nil
}
