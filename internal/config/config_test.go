package config_test

import (
	"bytes"
	"testing"

	"github.com/glasware/glas-core/internal/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Load(t *testing.T) {
	afs := afero.NewMemMapFs()
	writeTestConfig(t, afs)

	cfg, err := config.Load("./glas", config.OptAfs(afs))
	require.NoError(t, err)

	assert.Equal(t, "g^", cfg.Prefix())
	assert.True(t, cfg.Echo())

	require.Equal(t, 2, cfg.Aliases().Len())

	action, err := cfg.Aliases().Check("test bar")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = action(&buf)
	require.NoError(t, err)

	assert.Equal(t, "foo bar", buf.String())

	action, err = cfg.Aliases().Check("nested bar")
	require.NoError(t, err)

	buf.Reset()
	err = action(&buf)
	require.NoError(t, err)

	assert.Equal(t, "nested foo bar", buf.String())
}

func TestConfig_Options(t *testing.T) {
	afs := afero.NewMemMapFs()

	cfg, err := config.Load("./glas",
		config.OptAfs(afs),
		config.OptCommandPrefix("~"),
		config.OptNoEcho())
	require.NoError(t, err)

	assert.Equal(t, "~", cfg.Prefix())
	assert.False(t, cfg.Echo())
}

func writeTestConfig(t *testing.T, afs afero.Fs) {
	err := afs.Mkdir("./glas", 0765)
	require.NoError(t, err)

	hcl := `alias "test glob" {
	pattern = "test %s"
	template = "foo %s"
}`

	err = afero.WriteFile(afs, "./glas/test.hcl", []byte(hcl), 0765)
	require.NoError(t, err)

	err = afs.Mkdir("./glas/nested", 0765)
	require.NoError(t, err)

	hcl = `alias "test nested" {
	pattern = "nested %s"
	template = "nested foo %s"
}`

	err = afero.WriteFile(afs, "./glas/nested/test.hcl", []byte(hcl), 0765)
	require.NoError(t, err)
}
