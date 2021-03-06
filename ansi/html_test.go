package ansi

import (
	"testing"

	"github.com/glasware/glas-core/internal/test"
)

func TestHTML(t *testing.T) {
	testCases := []struct {
		d string
		e string
	}{
		{`no ansi codes are set`, `no ansi codes are set`},
		{`we pass a reset code\033[0m`, `we pass a reset code`},
		{`\033[40m\033[37mblack background with white text`, `black background with white text`},
		{`\033[40mblack background\033[32mgreen text \033[37mwhite text`, `black background<span style="color:#008000;">green text </span>white text`},
		{`\033[2J`, "$instruction$ERASESCREEN"},
		{`[2J[H7[1;24r8`, "$instruction$ERASESCREEN"},
	}

	for _, tc := range testCases {
		a := ReplaceCodes(tc.d)
		test.Equals(t, tc.e, a)
	}
}
