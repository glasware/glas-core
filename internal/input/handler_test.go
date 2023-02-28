package input_test

import (
	"context"
	"testing"

	"github.com/glasware/glas-core/internal/config"
	"github.com/glasware/glas-core/internal/connection"
	connMock "github.com/glasware/glas-core/internal/connection/mock"
	"github.com/glasware/glas-core/internal/input"
	"github.com/glasware/glas-core/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func surface(t *testing.T, ctrl *gomock.Controller, conn connection.Connection) *mock.MockSurface {
	t.Helper()
	surface := mock.NewMockSurface(ctrl)
	surface.EXPECT().Echo()
	surface.EXPECT().CommandPrefix().Return(config.DefaultCommandPrefix)
	surface.EXPECT().Connection().Return(conn)
	return surface
}

func TestInputHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := connMock.NewMockConnection(ctrl)
	conn.EXPECT().Connected().Return(false)
	h := input.New(surface(t, ctrl, conn))
	err := h.Input(context.Background(), "foobar")
	require.NoError(t, err)
}
