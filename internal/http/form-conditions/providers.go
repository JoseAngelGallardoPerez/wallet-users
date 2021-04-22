package form_conditions

func Providers() []interface{} {
	return []interface{}{
		NewDefaultUserClass,
		func(s *DefaultUserClass) *ConditionRegistry {
			service := NewConditionRegistry()
			if err := service.Register(s); err != nil {
				panic(err.Error())
			}

			return service
		},
	}
}
