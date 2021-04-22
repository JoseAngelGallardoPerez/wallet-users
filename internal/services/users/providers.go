package users

func Providers() []interface{} {
	return []interface{}{
		NewConfirmationCode,
		NewPermissionGroupsFiller,
		NewUserService,
		NewAttributeService,
		NewUserLoaderService,
		NewCompanyService,
	}
}
