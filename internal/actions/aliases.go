package actions

import "sync"

type (
	Aliases struct {
		mux  sync.RWMutex
		list []Alias
	}

	Alias interface {
		Match(str string) (bool, []any, error)
		Action(values ...any) Action
	}
)

func (a *Aliases) AddAlias(alias Alias) {
	a.mux.Lock()
	a.list = append(a.list, alias)
	a.mux.Unlock()
}

func (a Aliases) Check(str string) (Action, error) {
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
