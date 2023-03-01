package glob

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlob_Match(t *testing.T) {
	for input, tc := range map[string]struct {
		alias       Alias
		shouldMatch bool
		expected    []any
	}{
		"test foo": {
			alias: Alias{
				Pattern: "test %s",
			},
			shouldMatch: true,
			expected: []any{
				func() *string {
					s := "foo"
					return &s
				}(),
			},
		},
		"test": {
			alias: Alias{
				Pattern: "test %s",
			},
			shouldMatch: false,
		},
		"testing": {
			alias: Alias{
				Pattern: "test %s",
			},
			shouldMatch: false,
		},
		"l": {
			alias: Alias{
				Pattern: "ll",
			},
			shouldMatch: false,
		},
		"ll": {
			alias: Alias{
				Pattern: "ll",
			},
			shouldMatch: true,
			expected:    []any{},
		},
	} {
		t.Run(input, func(tt *testing.T) {
			matched, values, err := tc.alias.Match(input)
			require.NoError(tt, err)

			assert.Equal(tt, tc.shouldMatch, matched)
			if tc.shouldMatch {
				assert.Equal(tt, tc.expected, values)
			}
		})
	}
}

func TestGlob_Action(t *testing.T) {
	for input, tc := range map[string]struct {
		alias    Alias
		expected string
	}{
		"test bar": {
			alias: Alias{
				Pattern:  "test %s",
				Template: "foo %s",
			},
			expected: "foo bar\n",
		},
		"ll": {
			alias: Alias{
				Pattern: "ll",
				Template: `leave
out`,
			},
			expected: "leave\nout\n",
		},
		"bb": {
			alias: Alias{
				Pattern:  "bb",
				Template: "board boat;buy passage",
			},
			expected: "board boat\nbuy passage\n",
		},
	} {
		t.Run(input, func(tt *testing.T) {
			matched, values, err := tc.alias.Match(input)
			require.NoError(t, err)
			require.True(t, matched)

			action := tc.alias.Action(values...)
			var buf bytes.Buffer
			err = action(testWriter{&buf})
			require.NoError(tt, err)

			assert.Equal(tt, tc.expected, buf.String())
		})
	}
}

type testWriter struct {
	buf *bytes.Buffer
}

func (w testWriter) Write(p []byte) (int, error) {
	p = append(p, []byte("\n")...)
	return w.buf.Write(p)
}