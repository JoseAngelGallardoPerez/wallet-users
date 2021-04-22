package invites

func Providers() []interface{} {
	return []interface{}{
		NewCreator,
	}
}
