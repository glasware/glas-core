package regex

/*
type Regex struct {
	Name     string
	Pattern  string
	Template string

	pattern *regexp.Regexp
}

var _ actions.Alias = new(Regex)

func (a *Regex) Match(str string) (bool, error) {
	if a.pattern == nil {
		var err error
		a.pattern, err = regexp.Compile(a.Pattern)
		if err != nil {
			return false, err
		}
	}

	return a.pattern.MatchString(str), nil
}

func (a *Regex) Action(str string) Action {
	return func(output io.Writer, telnet io.Writer) error {
		template, err := regexp.Compile(a.Template)
		if err != nil {
			return err
		}

		return nil
	}
}
*/
