package actions

import (
	"strings"
	"sync"
)

type (
	Aliases struct {
		mux  sync.RWMutex
		list []Alias
	}

	Alias interface {
		Match(str string) (bool, []any, error)
		Action(values ...any) Action
		String() string
	}
)

func (a *Aliases) AddAlias(alias Alias) {
	a.mux.Lock()
	a.list = append(a.list, alias)
	a.mux.Unlock()
}

func (a *Aliases) Len() int {
	a.mux.RLock()
	defer a.mux.RUnlock()
	return len(a.list)
}

func (a *Aliases) Check(str string) (Action, error) {
	a.mux.RLock()
	defer a.mux.RUnlock()

	for _, alias := range a.list {
		ok, values, err := alias.Match(str)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		return alias.Action(values...), nil
	}

	return nil, nil
}

func (a *Aliases) List() string {
	a.mux.RLock()
	defer a.mux.RUnlock()

	aliases := make([]string, 0, len(a.list)+1)
	aliases = append(aliases, "Current aliases:")
	for _, alias := range a.list {
		aliases = append(aliases, alias.String())
	}

	return strings.Join(aliases, "\n")
}
