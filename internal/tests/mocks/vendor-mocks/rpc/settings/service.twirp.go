package settings

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/Confialink/wallet-settings/rpc/proto/settings"
)

type MockSettingsHandler struct {
	mock.Mock
}

func (m *MockSettingsHandler) List(ctx context.Context, req *settings.Request) (*settings.Response, error) {
	args := m.Called(ctx, req)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args[0].(*settings.Response), args.Error(1)
}

func (m *MockSettingsHandler) Get(ctx context.Context, req *settings.Request) (*settings.Response, error) {
	args := m.Called(ctx, req)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args[0].(*settings.Response), args.Error(1)
}
