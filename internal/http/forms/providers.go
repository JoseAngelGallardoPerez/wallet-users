package forms

func Providers() []interface{} {
	return []interface{}{
		NewFactory,
		NewUser,
	}
}
