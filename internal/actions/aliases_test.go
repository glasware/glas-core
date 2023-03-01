package actions_test

import (
	"bytes"
	"testing"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/actions/aliases/glob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAliases_AddAlias(t *testing.T) {
	var aliases actions.Aliases

	aliases.AddAlias(&glob.Alias{
		Pattern: "test %s",
	})

	assert.Equal(t, 1, aliases.Len())
}

func TestAliases_Check(t *testing.T) {
	var aliases actions.Aliases

	aliases.AddAlias(&glob.Alias{
		Pattern:  "test %s",
		Template: "foo %s",
	})

	action, err := aliases.Check("test bar")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = action(&buf)
	require.NoError(t, err)

	assert.Equal(t, "foo bar", buf.String())
}
