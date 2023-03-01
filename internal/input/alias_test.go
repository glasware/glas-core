package input_test

import (
	"context"
	"testing"

	"github.com/glasware/glas-core/internal/actions"
	"github.com/glasware/glas-core/internal/actions/aliases/glob"
	"github.com/glasware/glas-core/internal/config"
	"github.com/glasware/glas-core/internal/connection/mock"
	"github.com/glasware/glas-core/internal/input"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestInputHandler_handleAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mock.NewMockConnection(ctrl)
	conn.EXPECT().Connected().Return(true)
	surface := surface(t, ctrl, conn)

	h := input.New(surface)

	alias := glob.Alias{
		Pattern:  "test %s",
		Template: "foo %s",
	}
	m := actions.Aliases{}
	m.AddAlias(&alias)

	surface.EXPECT().
		Aliases().
		Return(&m)

	conn.EXPECT().
		Write(gomock.Eq([]byte("foobar"))).
		Return(len([]byte("foobar")), nil)

	err := h.Input(context.Background(), "foobar")
	require.NoError(t, err)

	surface.EXPECT().Echo()
	surface.EXPECT().CommandPrefix().Return(config.DefaultCommandPrefix)
	surface.EXPECT().Connection().Return(conn)
	conn.EXPECT().Connected().Return(true)

	surface.EXPECT().
		Aliases().
		Return(&m)

	surface.EXPECT().Echo()
	surface.EXPECT().Connection().Return(conn)

	conn.EXPECT().
		Write(gomock.Eq([]byte("foo baz"))).
		Return(len([]byte("foo baz")), nil)

	err = h.Input(context.Background(), "test baz")
	require.NoError(t, err)
}
