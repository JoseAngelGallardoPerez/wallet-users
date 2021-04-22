package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/Confialink/wallet-users/internal/http/handlers"
	"github.com/Confialink/wallet-users/internal/services/permissions"

	"github.com/Confialink/wallet-users/internal/http/responses"
)

type PermissionsMiddleware struct {
	permissionsService *permissions.Permissions
	responseService    responses.ResponseHandler
}

func NewPermissionsMiddleware(
	permissionsService *permissions.Permissions,
	responseService responses.ResponseHandler,
) *PermissionsMiddleware {
	return &PermissionsMiddleware{permissionsService, responseService}
}

func (m *PermissionsMiddleware) CanViewProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		requestedUser := handlers.GetRequestedUser(ctx)
		if requestedUser == nil {
			m.responseService.NotFound(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanViewProfile(currentUser.UID, requestedUser); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanViewClientProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanViewUserProfile(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}
func (m *PermissionsMiddleware) CanViewAdminProfile() gin.HandlerFunc {
	// at current time we have one permission for admin and client profiles
	return m.CanViewClientProfile()
}

func (m *PermissionsMiddleware) CanViewShortUserProfiles() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanViewShortUserProfiles(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanUpdateProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		requestedUser := handlers.GetRequestedUser(ctx)
		if requestedUser == nil {
			m.responseService.NotFound(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanUpdateProfile(currentUser.UID, requestedUser); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanViewSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanViewSettings(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanModifySettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanModifySettings(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanCreateSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanCreateSettings(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}

func (m *PermissionsMiddleware) CanRemoveSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := handlers.GetCurrentUser(ctx)
		if currentUser == nil {
			m.responseService.Forbidden(ctx)
			return
		}

		if hasPerm := m.permissionsService.CanRemoveSettings(currentUser.UID); !hasPerm {
			m.responseService.Forbidden(ctx)
			return
		}
	}
}
