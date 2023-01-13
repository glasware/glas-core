package aliases

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/glasware/glas-core/internal/actions"
)

type Glob struct {
	Name     string
	Pattern  string
	Template string

	calcArgsOnce sync.Once
	argTypes     []reflect.Kind
}

var _ actions.Alias = new(Glob)

func (a *Glob) Match(str string) (bool, []any, error) {
	values, err := a.valuePtrs()
	if err != nil {
		return false, nil, err
	}

	n, err := fmt.Sscanf(str, a.Pattern, values...)
	if err != nil {
		if err.Error() == "input does not match format" {
			return false, nil, nil
		}

		return false, nil, err
	}

	if n != len(values) {
		return false, nil, nil
	}

	return true, values, nil
}

func (a *Glob) Action(in ...any) actions.Action {
	template := fmt.Sprintf(a.Template, a.values(in...)...)
	commands := strings.Split(template, "\n")
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

func (a *Glob) String() string {
	return fmt.Sprintf("%s | %s = %s", a.Name, a.Pattern, a.Template)
}

func (a *Glob) values(in ...any) []any {
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

func (a *Glob) valuePtrs() ([]any, error) {
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

func (a *Glob) scan() error {
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
