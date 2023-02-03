package input_test

import (
	"testing"

	"github.com/glasware/glas-core/internal/input"
	"github.com/golang/mock/gomock"
	"github.com/glasware/glas-core/internal/mock"
)

func TestInputHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	surface := mock.NewMockSurface(ctrl)

	h := input.New(surface)
	h = h
}
