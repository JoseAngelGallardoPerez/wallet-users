package auth

func Providers() []interface{} {
	return []interface{}{
		NewSignUpResponse,
	}
}
