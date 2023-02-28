package aliases

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlob_Match(t *testing.T) {
	for input, tc := range map[string]struct {
		alias       Glob
		shouldMatch bool
		expected    []any
	}{
		"test foo": {
			alias: Glob{
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
			alias: Glob{
				Pattern: "test %s",
			},
			shouldMatch: false,
		},
		"testing": {
			alias: Glob{
				Pattern: "test %s",
			},
			shouldMatch: false,
		},
		"l": {
			alias: Glob{
				Pattern: "ll",
			},
			shouldMatch: false,
		},
		"ll": {
			alias: Glob{
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
		alias    Glob
		expected string
	}{
		"test bar": {
			alias: Glob{
				Pattern:  "test %s",
				Template: "foo %s",
			},
			expected: "foo bar\n",
		},
		"ll": {
			alias: Glob{
				Pattern: "ll",
				Template: `leave
out`,
			},
			expected: "leave\nout\n",
		},
		"bb": {
			alias: Glob{
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
