package repositories

func Providers() []interface{} {
	return []interface{}{
		NewRepository,
		NewAccesLogRepository,
		NewBlockedIpsRepository,
		NewConfirmationCodeRepository,
		NewFailAuthAttemptRepository,
		NewUserGroupsRepository,
		NewSecurityQuestionRepository,
		NewSecurityQuestionsAnswerRepository,
		NewTokenRepository,
		NewUsersRepository,
		NewVerificationRepository,
		NewVerificationFilesRepository,
		NewInvitesRepository,
		NewFormRepository,
		NewAddressRepository,
		NewUserAttributeValueRepository,
		NewAttributeRepository,
		NewCompanyRepository,
	}
}
