package messagebroker

func Providers() []interface{} {
	return []interface{}{
		NewNats,
	}
}
