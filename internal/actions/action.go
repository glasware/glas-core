package actions

import "io"

type Action = func(writer io.Writer) error
