package glob

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/glasware/glas-core/internal/actions"
)

type Alias struct {
	Name     string
	Pattern  string
	Template string

	calcArgsOnce sync.Once
	argTypes     []reflect.Kind
}

var _ actions.Alias = new(Alias)

func (a *Alias) Match(str string) (bool, []any, error) {
	values, err := a.valuePtrs()
	if err != nil {
		return false, nil, fmt.Errorf("valuePtrs: %w", err)
	}

	n, err := fmt.Sscanf(str, a.Pattern, values...)
	if err != nil {
		if err := shouldError(err); err != nil {
			return false, nil, fmt.Errorf("fmt.Sscanf: %w", err)
		}

		return false, nil, nil
	}

	if n != len(values) {
		return false, nil, nil
	}

	return true, values, nil
}

func (a *Alias) Action(in ...any) actions.Action {
	template := fmt.Sprintf(a.Template, a.values(in...)...)
	commands := strings.FieldsFunc(template, func(c rune) bool {
		return c == '\n' || c == ';'
	})

	return func(writer io.Writer) error {
		for _, command := range commands {
			cmd := []byte(command)
			if _, err := writer.Write(cmd); err != nil {
				return err
			}
		}

		return nil
	}
}

func (a *Alias) String() string {
	return fmt.Sprintf("%s | %s = %s", a.Name, a.Pattern, a.Template)
}

func (a *Alias) values(in ...any) []any {
	values := make([]any, 0, len(in))
	for _, v := range in {
		switch nv := v.(type) {
		case *string:
			values = append(values, *nv)
		case *int:
			values = append(values, *nv)
		case *bool:
			values = append(values, *nv)
		case *interface{}:
			values = append(values, *nv)
		}
	}
	return values
}

func (a *Alias) valuePtrs() ([]any, error) {
	var err error
	a.calcArgsOnce.Do(func() {
		err = a.scan()
	})
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, 0, len(a.argTypes))
	for _, at := range a.argTypes {
		switch at {
		case reflect.String:
			var arg string
			values = append(values, &arg)
		case reflect.Int:
			var arg int
			values = append(values, &arg)
		case reflect.Bool:
			var arg bool
			values = append(values, &arg)
		case reflect.Interface:
			var arg interface{}
			values = append(values, &arg)
		}
	}

	return values, nil
}

func (a *Alias) scan() error {
	scanner := bufio.NewScanner(strings.NewReader(a.Pattern))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		switch scanner.Text() {
		case "%s":
			a.argTypes = append(a.argTypes, reflect.String)
		case "%d":
			a.argTypes = append(a.argTypes, reflect.Int)
		case "%t":
			a.argTypes = append(a.argTypes, reflect.Bool)
		case "%v":
			a.argTypes = append(a.argTypes, reflect.Interface)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}