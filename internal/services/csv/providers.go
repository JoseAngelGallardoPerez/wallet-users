package csv

func Providers() []interface{} {
	return []interface{}{
		NewAdminProfiles,
		NewUsers,
	}
}
