package notification

import (
	"context"

	"github.com/Confialink/wallet-notifications/rpc/proto/notifications"
	"github.com/stretchr/testify/mock"
)

type MockNotificationHandler struct {
	mock.Mock
}

func (m *MockNotificationHandler) Dispatch(ctx context.Context, in *notifications.Request) (*notifications.Response, error) {
	args := m.Called(ctx, in)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args[0].(*notifications.Response), args.Error(1)
}

func (m *MockNotificationHandler) GetSettings(ctx context.Context, in *notifications.SettingsRequest) (*notifications.SettingsResponse, error) {
	args := m.Called(ctx, in)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args[0].(*notifications.SettingsResponse), args.Error(1)
}

func (m *MockNotificationHandler) GetUserSettings(ctx context.Context, in *notifications.UserSettingsRequest) (*notifications.UserSettingsResponse, error) {
	args := m.Called(ctx, in)
	if args[0] == nil {
		return nil, args.Error(1)
	}

	return args[0].(*notifications.UserSettingsResponse), args.Error(1)
}
