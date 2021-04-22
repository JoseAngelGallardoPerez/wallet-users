package usersserver

func Providers() []interface{} {
	return []interface{}{
		NewUsersServer,
	}
}
