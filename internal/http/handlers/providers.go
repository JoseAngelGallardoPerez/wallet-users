package handlers

func Providers() []interface{} {
	return []interface{}{
		NewAuthService,
		NewBlockedIpsService,
		NewUserGroupsService,
		NewHandlerParams,
		NewSecurityQuestionService,
		NewSecurityQuestionAnswerService,
		NewUsersService,
		NewVerificationHandler,
		NewStaffsService,
		NewInvitesHandler,
	}
}
