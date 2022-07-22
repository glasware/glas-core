package aliases

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlob_Match(t *testing.T) {
	alias := &Glob{
		Pattern: "test %s",
	}

	matched, values, err := alias.Match("test foo")
	require.NoError(t, err)

	require.True(t, matched)
	assert.Equal(t, []any{
		func() *string {
			s := "foo"
			return &s
		}(),
	}, values)
}

func TestGlob_Action(t *testing.T) {
	alias := &Glob{
		Pattern:  "test %s",
		Template: "foo %s",
	}

	matched, values, err := alias.Match("test bar")
	require.NoError(t, err)
	require.True(t, matched)

	action := alias.Action(values...)
	var buf bytes.Buffer
	err = action(&buf)
	require.NoError(t, err)

	assert.Equal(t, "foo bar", buf.String())
}
