package auth

func Providers() []interface{} {
	return []interface{}{
		BlockerFactory,
		TokenServiceFactory,
		NewAutologoutTTLResolver,
		NewFixedValueTTLResolver,
		TemporaryTokensFactory,
		NewAuth,
	}
}
